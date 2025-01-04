package startrpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/chat/pkg/common/kdisc"
	disetcd "github.com/openimsdk/chat/pkg/common/kdisc/etcd"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/runtimeenv"

	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/network"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Start rpc server.
func Start[T any](ctx context.Context, discovery *config.Discovery, listenIP,
	registerIP string, rpcPorts []int, index int, rpcRegisterName string, share *config.Share, config T,
	watchConfigNames []string, watchServiceNames []string,
	rpcFn func(ctx context.Context, config T, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {

	runtimeEnv := runtimeenv.PrintRuntimeEnvironment()

	rpcPort, err := datautil.GetElemByIndex(rpcPorts, index)
	if err != nil {
		return err
	}
	log.CInfo(ctx, "RPC server is initializing", " runtimeEnv ", runtimeEnv, "rpcRegisterName", rpcRegisterName, "rpcPort", rpcPort)
	rpcTcpAddr := net.JoinHostPort(network.GetListenIP(listenIP), strconv.Itoa(rpcPort))
	listener, err := net.Listen(
		"tcp",
		rpcTcpAddr,
	)
	if err != nil {
		return errs.WrapMsg(err, "listen err", "rpcTcpAddr", rpcTcpAddr)
	}

	defer listener.Close()
	client, err := kdisc.NewDiscoveryRegister(discovery, runtimeEnv, watchServiceNames)
	if err != nil {
		return err
	}
	defer client.Close()
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	registerIP, err = network.GetRpcRegisterIP(registerIP)
	if err != nil {
		return err
	}

	options = append(options, mw.GrpcServer())
	srv := grpc.NewServer(options...)
	once := sync.Once{}
	defer func() {
		once.Do(srv.GracefulStop)
	}()

	err = rpcFn(ctx, config, client, srv)
	if err != nil {
		return err
	}

	if err := client.Register(rpcRegisterName, registerIP, rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	}

	var (
		netDone = make(chan struct{}, 2)
		netErr  error
	)

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			netErr = errs.WrapMsg(err, "rpc start err: ", rpcTcpAddr)
			netDone <- struct{}{}
		}
	}()
	if discovery.Enable == kdisc.ETCDCONST {
		cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(), watchConfigNames)
		cm.Watch(ctx)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := gracefulStopWithCtx(ctx, srv.GracefulStop); err != nil {
			return err
		}
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return nil
	case <-netDone:
		close(netDone)
		return netErr
	}
}

func gracefulStopWithCtx(ctx context.Context, f func()) error {
	done := make(chan struct{}, 1)
	go func() {
		f()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return errs.New("timeout, ctx graceful stop")
	case <-done:
		return nil
	}
}
