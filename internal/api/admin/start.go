package admin

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
	SetAdminRoute(engine, adminApi, mwApi)
	return engine.Run(fmt.Sprintf(":%d", apiPort))
}

func SetAdminRoute(router gin.IRouter, admin *Api, mw *chatmw.MW) {

	adminRouterGroup := router.Group("/account")
	adminRouterGroup.POST("/login", admin.AdminLogin)                                   // Login
	adminRouterGroup.POST("/update", mw.CheckAdmin, admin.AdminUpdateInfo)              // Modify information
	adminRouterGroup.POST("/info", mw.CheckAdmin, admin.AdminInfo)                      // Get information
	adminRouterGroup.POST("/change_password", mw.CheckAdmin, admin.ChangeAdminPassword) // Change admin account's password
	adminRouterGroup.POST("/add_admin", mw.CheckAdmin, admin.AddAdminAccount)           // Add admin account
	adminRouterGroup.POST("/add_user", mw.CheckAdmin, admin.AddUserAccount)             // Add user account
	adminRouterGroup.POST("/del_admin", mw.CheckAdmin, admin.DelAdminAccount)           // Delete admin
	adminRouterGroup.POST("/search", mw.CheckAdmin, admin.SearchAdminAccount)           // Get admin list
	//account.POST("/add_notification_account")

	importGroup := router.Group("/user/import")
	importGroup.POST("/json", mw.CheckAdmin, admin.ImportUserByJson)
	importGroup.POST("/xlsx", mw.CheckAdmin, admin.ImportUserByXlsx)
	importGroup.GET("/xlsx", admin.BatchImportTemplate)

	allowRegisterGroup := router.Group("/user/allow_register", mw.CheckAdmin)
	allowRegisterGroup.POST("/get", admin.GetAllowRegister)
	allowRegisterGroup.POST("/set", admin.SetAllowRegister)

	defaultRouter := router.Group("/default", mw.CheckAdmin)
	defaultUserRouter := defaultRouter.Group("/user")
	defaultUserRouter.POST("/add", admin.AddDefaultFriend)       // Add default friend at registration
	defaultUserRouter.POST("/del", admin.DelDefaultFriend)       // Delete default friend at registration
	defaultUserRouter.POST("/find", admin.FindDefaultFriend)     // Default friend list
	defaultUserRouter.POST("/search", admin.SearchDefaultFriend) // Search default friend list at registration
	defaultGroupRouter := defaultRouter.Group("/group")
	defaultGroupRouter.POST("/add", admin.AddDefaultGroup)       // Add default group at registration
	defaultGroupRouter.POST("/del", admin.DelDefaultGroup)       // Delete default group at registration
	defaultGroupRouter.POST("/find", admin.FindDefaultGroup)     // Get default group list at registration
	defaultGroupRouter.POST("/search", admin.SearchDefaultGroup) // Search default group list at registration

	invitationCodeRouter := router.Group("/invitation_code", mw.CheckAdmin)
	invitationCodeRouter.POST("/add", admin.AddInvitationCode)       // Add invitation code
	invitationCodeRouter.POST("/gen", admin.GenInvitationCode)       // Generate invitation code
	invitationCodeRouter.POST("/del", admin.DelInvitationCode)       // Delete invitation code
	invitationCodeRouter.POST("/search", admin.SearchInvitationCode) // Search invitation code

	forbiddenRouter := router.Group("/forbidden", mw.CheckAdmin)
	ipForbiddenRouter := forbiddenRouter.Group("/ip")
	ipForbiddenRouter.POST("/add", admin.AddIPForbidden)       // Add forbidden IP for registration/login
	ipForbiddenRouter.POST("/del", admin.DelIPForbidden)       // Delete forbidden IP for registration/login
	ipForbiddenRouter.POST("/search", admin.SearchIPForbidden) // Search forbidden IPs for registration/login
	userForbiddenRouter := forbiddenRouter.Group("/user")
	userForbiddenRouter.POST("/add", admin.AddUserIPLimitLogin)       // Add limit for user login on specific IP
	userForbiddenRouter.POST("/del", admin.DelUserIPLimitLogin)       // Delete user limit on specific IP for login
	userForbiddenRouter.POST("/search", admin.SearchUserIPLimitLogin) // Search limit for user login on specific IP

	appletRouterGroup := router.Group("/applet", mw.CheckAdmin)
	appletRouterGroup.POST("/add", admin.AddApplet)       // Add applet
	appletRouterGroup.POST("/del", admin.DelApplet)       // Delete applet
	appletRouterGroup.POST("/update", admin.UpdateApplet) // Modify applet
	appletRouterGroup.POST("/search", admin.SearchApplet) // Search applet

	blockRouter := router.Group("/block", mw.CheckAdmin)
	blockRouter.POST("/add", admin.BlockUser)          // Block user
	blockRouter.POST("/del", admin.UnblockUser)        // Unblock user
	blockRouter.POST("/search", admin.SearchBlockUser) // Search blocked users

	userRouter := router.Group("/user", mw.CheckAdmin)
	userRouter.POST("/password/reset", admin.ResetUserPassword) // Reset user password

	initGroup := router.Group("/client_config", mw.CheckAdmin)
	initGroup.POST("/get", admin.GetClientConfig) // Get client initialization configuration
	initGroup.POST("/set", admin.SetClientConfig) // Set client initialization configuration
	initGroup.POST("/del", admin.DelClientConfig) // Delete client initialization configuration

	statistic := router.Group("/statistic", mw.CheckAdmin)
	statistic.POST("/new_user_count", admin.NewUserCount)
	statistic.POST("/login_user_count", admin.LoginUserCount)

	applicationGroup := router.Group("application")
	applicationGroup.POST("/add_version", mw.CheckAdmin, admin.AddApplicationVersion)
	applicationGroup.POST("/update_version", mw.CheckAdmin, admin.UpdateApplicationVersion)
	applicationGroup.POST("/delete_version", mw.CheckAdmin, admin.DeleteApplicationVersion)
	applicationGroup.POST("/latest_version", admin.LatestApplicationVersion)
	applicationGroup.POST("/page_versions", admin.PageApplicationVersion)
}
