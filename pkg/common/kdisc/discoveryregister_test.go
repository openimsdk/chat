package kdisc

import (
	"testing"

	"github.com/openimsdk/chat/pkg/common/config"
	toolsetcd "github.com/openimsdk/tools/discovery/etcd"
)

func TestNewDiscoveryRegister_selectsEtcd_whenConfiguredInKubernetes(t *testing.T) {
	// Given
	discoveryConfig := &config.Discovery{
		Enable: ETCDCONST,
		Etcd: config.Etcd{
			RootDirectory: "openim",
			Address:       []string{"127.0.0.1:2379"},
		},
	}

	// When
	registry, err := NewDiscoveryRegister(discoveryConfig, KUBERNETESCONST, nil)

	// Then
	if err != nil {
		t.Fatalf("NewDiscoveryRegister() error = %v", err)
	}
	defer registry.Close()
	if _, ok := registry.(*toolsetcd.SvcDiscoveryRegistryImpl); !ok {
		t.Fatalf("NewDiscoveryRegister() type = %T, want *etcd.SvcDiscoveryRegistryImpl", registry)
	}
}
