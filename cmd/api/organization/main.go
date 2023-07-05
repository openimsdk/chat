package main

func main() {
	//rand.Seed(time.Now().UnixNano())
	//err := config.InitConfig()
	//if err != nil {
	//	panic(err)
	//}
	//if err := log.InitFromConfig("organization.log", "organization-api", *config.Config.Log.RemainLogLevel, *config.Config.Log.IsStdout, *config.Config.Log.IsJson, *config.Config.Log.StorageLocation, *config.Config.Log.RemainRotationCount); err != nil {
	//	panic(err)
	//}
	//zk, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
	//	openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.UserName,
	//		config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	//if err != nil {
	//	panic(err)
	//}
	//zk.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	//engine := gin.Default()
	//engine.Use(mw.CorsHandler(), mw.GinParseOperationID())
	//api.NewOrganizationRoute(engine, zk)
	//
	//
	//defaultPorts := config.Config.OrganizationApi.GinPort
	//ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd")
	//flag.Parse()
	//address := net.JoinHostPort(config.Config.ChatApi.ListenIP, strconv.Itoa(*ginPort))
	//if err := engine.Run(address); err != nil {
	//	panic(err)
	//}
}
