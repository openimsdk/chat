package organization

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/proto/organization"
)

type OrgClient struct {
	client organization.OrganizationClient
}

func NewOrgClient(discov discoveryregistry.SvcDiscoveryRegistry) *OrgClient {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImOrganizationName)
	if err != nil {
		panic(err)
	}
	return &OrgClient{
		client: organization.NewOrganizationClient(conn),
	}
}
