package bot

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
	"github.com/openimsdk/chat/pkg/common/kdisc"
	disetcd "github.com/openimsdk/chat/pkg/common/kdisc/etcd"
	adminclient "github.com/openimsdk/chat/pkg/protocol/admin"
	botclient "github.com/openimsdk/chat/pkg/protocol/bot"
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
	ApiConfig config.APIBot
	Discovery config.Discovery
	Share     config.Share
	Redis     config.Redis

	RuntimeEnv string
}

func Start(ctx context.Context, index int, cfg *Config) error {
	cfg.RuntimeEnv = runtimeenv.PrintRuntimeEnvironment()
	apiPort, err := datautil.GetElemByIndex(cfg.ApiConfig.Api.Ports, index)
	if err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(&cfg.Discovery, cfg.RuntimeEnv, nil)
	if err != nil {
		return err
	}

	botConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Bot, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	adminConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Admin, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	adminClient := adminclient.NewAdminClient(adminConn)
	botClient := botclient.NewBotClient(botConn)
	base := util.Api{
		ImUserID:        cfg.Share.OpenIM.AdminUserID,
		ProxyHeader:     cfg.Share.ProxyHeader,
		ChatAdminUserID: cfg.Share.ChatAdmin[0],
	}
	botApi := New(botClient, &base)
	mwApi := chatmw.New(adminClient)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	SetBotRoute(engine, botApi, mwApi)

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
				config.ChatAPIBotCfgFileName,
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

func SetBotRoute(router gin.IRouter, bot *Api, mw *chatmw.MW) {
	account := router.Group("/agent")
	account.POST("/create", mw.CheckAdmin, bot.CreateAgent)
	account.POST("/delete", mw.CheckAdmin, bot.DeleteAgent)
	account.POST("/update", mw.CheckAdmin, bot.UpdateAgent)
	account.POST("/page", mw.CheckToken, bot.PageFindAgent)

	imwebhook := router.Group("/im_callback")
	imwebhook.POST("/callbackAfterSendSingleMsgCommand", bot.AfterSendSingleMsg)
	imwebhook.POST("/callbackAfterSendGroupMsgCommand", bot.AfterSendGroupMsg)
}
