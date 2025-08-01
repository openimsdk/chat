package chat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	chatmw "github.com/openimsdk/chat/internal/api/mw"
	"github.com/openimsdk/chat/internal/api/util"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/common/kdisc"
	disetcd "github.com/openimsdk/chat/pkg/common/kdisc/etcd"
	adminclient "github.com/openimsdk/chat/pkg/protocol/admin"
	chatclient "github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	ApiConfig config.API
	Discovery config.Discovery
	Share     config.Share
	Redis     config.Redis

	RuntimeEnv string
}

func Start(ctx context.Context, index int, cfg *Config) error {
	cfg.RuntimeEnv = runtimeenv.PrintRuntimeEnvironment()

	if len(cfg.Share.ChatAdmin) == 0 {
		return errs.New("share chat admin not configured")
	}
	apiPort, err := datautil.GetElemByIndex(cfg.ApiConfig.Api.Ports, index)
	if err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(&cfg.Discovery, cfg.RuntimeEnv, nil)
	if err != nil {
		return err
	}

	chatConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Chat, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	adminConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Admin, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	chatClient := chatclient.NewChatClient(chatConn)
	adminClient := adminclient.NewAdminClient(adminConn)
	im := imapi.New(cfg.Share.OpenIM.ApiURL, cfg.Share.OpenIM.Secret, cfg.Share.OpenIM.AdminUserID)
	base := util.Api{
		ImUserID:        cfg.Share.OpenIM.AdminUserID,
		ProxyHeader:     cfg.Share.ProxyHeader,
		ChatAdminUserID: cfg.Share.ChatAdmin[0],
	}
	adminApi := New(chatClient, adminClient, im, &base)
	mwApi := chatmw.New(adminClient)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	SetChatRoute(engine, adminApi, mwApi)

	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)
	server := http.Server{Addr: fmt.Sprintf(":%d", apiPort), Handler: engine}
	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			netErr = errs.WrapMsg(err, fmt.Sprintf("api start err: %s", server.Addr))
			netDone <- struct{}{}
		}
	}()
	if cfg.Discovery.Enable == kdisc.ETCDCONST {
		cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(),
			[]string{
				config.ChatAPIChatCfgFileName,
				config.DiscoveryConfigFileName,
				config.ShareFileName,
				config.LogConfigFileName,
			},
		)
		cm.Watch(ctx)
	}
	shutdown := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			return errs.WrapMsg(err, "shutdown err")
		}
		return nil
	}
	disetcd.RegisterShutDown(shutdown)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
		if err := shutdown(); err != nil {
			return err
		}
	case <-netDone:
		close(netDone)
		return netErr
	}
	return nil
}

func SetChatRoute(router gin.IRouter, chat *Api, mw *chatmw.MW) {
	account := router.Group("/account")
	account.POST("/code/send", chat.SendVerifyCode)                      // Send verification code
	account.POST("/code/verify", chat.VerifyCode)                        // Verify the verification code
	account.POST("/register", mw.CheckAdminOrNil, chat.RegisterUser)     // Register
	account.POST("/login", chat.Login)                                   // Login
	account.POST("/password/reset", chat.ResetPassword)                  // Forgot password
	account.POST("/password/change", mw.CheckToken, chat.ChangePassword) // Change password

	user := router.Group("/user", mw.CheckToken)
	user.POST("/update", chat.UpdateUserInfo)                 // Edit personal information
	user.POST("/find/public", chat.FindUserPublicInfo)        // Get user's public information
	user.POST("/find/full", chat.FindUserFullInfo)            // Get all information of the user
	user.POST("/search/full", chat.SearchUserFullInfo)        // Search user's public information
	user.POST("/search/public", chat.SearchUserPublicInfo)    // Search all information of the user
	user.POST("/rtc/get_token", chat.GetTokenForVideoMeeting) // Get token for video meeting for the user

	router.POST("/friend/search", mw.CheckToken, chat.SearchFriend)

	router.Group("/applet").POST("/find", mw.CheckToken, chat.FindApplet) // Applet list

	router.Group("/client_config").POST("/get", chat.GetClientConfig) // Get client initialization configuration

	applicationGroup := router.Group("application")
	applicationGroup.POST("/latest_version", chat.LatestApplicationVersion)
	applicationGroup.POST("/page_versions", chat.PageApplicationVersion)

	router.Group("/callback").POST("/open_im", chat.OpenIMCallback) // Callback
}
