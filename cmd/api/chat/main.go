package main

import (
	"flag"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw"
	"github.com/OpenIMSDK/chat/internal/api"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/openKeeper"
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
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	if err := log.InitFromConfig("chat.log", "chat-api", config.Config.Log.RemainLogLevel, config.Config.Log.IsStdout, config.Config.Log.IsJson, config.Config.Log.StorageLocation, config.Config.Log.RemainRotationCount); err != nil {
		panic(err)
	}
	zk, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.UserName,
			config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10))
	if err != nil {
		panic(err)
	}
	zk.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	engine := gin.Default()
	engine.Use(mw.CorsHandler(), mw.GinParseOperationID())
	api.NewChatRoute(engine, zk)
	defaultPorts := config.Config.ChatApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd")
	flag.Parse()
	address := net.JoinHostPort(config.Config.ChatApi.ListenIP, strconv.Itoa(*ginPort))
	if err := engine.Run(address); err != nil {
		panic(err)
	}
}
