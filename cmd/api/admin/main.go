package main

import (
	"flag"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw"
	openKeeper "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/chat/internal/api"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := log.InitFromConfig("chat.log", "admin-api", config.Config.Log.RemainLogLevel, config.Config.Log.IsStdout, config.Config.Log.IsJson, config.Config.Log.StorageLocation, config.Config.Log.RemainRotationCount); err != nil {
		panic(err)
	}
	zk, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.UserName,
			config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		panic(err)
	}
	zk.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	engine := gin.Default()
	engine.Use(mw.CorsHandler(), mw.GinParseOperationID())
	api.NewAdminRoute(engine, zk)
	defaultPorts := config.Config.AdminApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd")
	flag.Parse()
	address := net.JoinHostPort(config.Config.AdminApi.ListenIP, strconv.Itoa(*ginPort))
	if err := engine.Run(address); err != nil {
		panic(err)
	}
}
