package kdisc

import (
	"time"

	"github.com/openimsdk/chat/pkg/common/config"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/discovery/kubernetes"
	"github.com/openimsdk/tools/errs"
)

const (
	ETCDCONST       = "etcd"
	KUBERNETESCONST = "kubernetes"
	DIRECTCONST     = "direct"
)

// NewDiscoveryRegister creates a new service discovery and registry client based on the provided environment type.
func NewDiscoveryRegister(discovery *config.Discovery, runtimeEnv string, watchNames []string) (discovery.SvcDiscoveryRegistry, error) {
	if runtimeEnv == KUBERNETESCONST {
		return kubernetes.NewKubernetesConnManager(discovery.Kubernetes.Namespace)
	}

	switch discovery.Enable {
	case ETCDCONST:
		return etcd.NewSvcDiscoveryRegistry(
			discovery.Etcd.RootDirectory,
			discovery.Etcd.Address,
			watchNames,
			etcd.WithDialTimeout(10*time.Second),
			etcd.WithMaxCallSendMsgSize(20*1024*1024),
			etcd.WithUsernameAndPassword(discovery.Etcd.Username, discovery.Etcd.Password))
	default:
		return nil, errs.New("unsupported discovery type", "type", discovery.Enable).Wrap()
	}
}
