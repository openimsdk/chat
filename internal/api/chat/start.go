package chat

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	chatmw "github.com/openimsdk/chat/internal/api/mw"
	"github.com/openimsdk/chat/internal/api/util"
	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/common/kdisc"
	adminclient "github.com/openimsdk/chat/pkg/protocol/admin"
	chatclient "github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	ApiConfig config.API
	Discovery config.Discovery
	Share     config.Share
}

func Start(ctx context.Context, index int, config *Config) error {
	if len(config.Share.ChatAdmin) == 0 {
		return errs.New("share chat admin not configured")
	}
	apiPort, err := datautil.GetElemByIndex(config.ApiConfig.Api.Ports, index)
	if err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(&config.Discovery)
	if err != nil {
		return err
	}

	chatConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Chat, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	adminConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Admin, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		return err
	}
	chatClient := chatclient.NewChatClient(chatConn)
	adminClient := adminclient.NewAdminClient(adminConn)
	im := imapi.New(config.Share.OpenIM.ApiURL, config.Share.OpenIM.Secret, config.Share.OpenIM.AdminUserID)
	base := util.Api{
		ImUserID:        config.Share.OpenIM.AdminUserID,
		ProxyHeader:     config.Share.ProxyHeader,
		ChatAdminUserID: config.Share.ChatAdmin[0],
	}
	adminApi := New(chatClient, adminClient, im, &base)
	mwApi := chatmw.New(adminClient)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	SetChatRoute(engine, adminApi, mwApi)
	return engine.Run(fmt.Sprintf(":%d", apiPort))
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
