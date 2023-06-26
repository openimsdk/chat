package main

import (
	"flag"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw"
	openKeeper "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/chat/internal/api"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	if err := log.InitFromConfig("organization.log", "organization-api", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount); err != nil {
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
	api.NewOrganizationRoute(engine, zk)

	client := http.Client{Timeout: time.Second * 60}
	engine.NoRoute(func(c *gin.Context) {
		c.Abort()
		reqUrl := config.Config.BusinessServerAddress + c.Request.RequestURI
		req, err := http.NewRequest(c.Request.Method, reqUrl, c.Request.Body)
		if err != nil {
			c.String(500, fmt.Sprintf("BusinessServer [%s] %s req=> %s", c.Request.Method, reqUrl, err))
			return
		}
		req.Header = c.Request.Header
		resp, err := client.Do(req)
		if err != nil {
			c.String(500, fmt.Sprintf("BusinessServer [%s] %s do=> %s", c.Request.Method, reqUrl, err))
			return
		}
		for name, values := range resp.Header {
			c.Writer.Header().Del(name)
			for _, val := range values {
				c.Writer.Header().Add(name, val)
			}
		}
		c.Status(resp.StatusCode)
		c.Writer.WriteHeaderNow()
		io.Copy(c.Writer, resp.Body)
	})

	defaultPorts := config.Config.OrganizationApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd")
	flag.Parse()
	address := net.JoinHostPort(config.Config.ChatApi.ListenIP, strconv.Itoa(*ginPort))
	if err := engine.Run(address); err != nil {
		panic(err)
	}
}
