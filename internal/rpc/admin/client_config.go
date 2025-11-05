package admin

import (
	"context"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) GetClientConfig(ctx context.Context, req *admin.GetClientConfigReq) (*admin.GetClientConfigResp, error) {
	conf, err := o.Database.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &admin.GetClientConfigResp{Config: conf}, nil
}

func (o *adminServer) SetClientConfig(ctx context.Context, req *admin.SetClientConfigReq) (*admin.SetClientConfigResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.Config) == 0 {
		return nil, errs.ErrArgs.WrapMsg("update config empty")
	}
	if err := o.Database.SetConfig(ctx, req.Config); err != nil {
		return nil, err
	}
	return &admin.SetClientConfigResp{}, nil
}

func (o *adminServer) DelClientConfig(ctx context.Context, req *admin.DelClientConfigReq) (*admin.DelClientConfigResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := o.Database.DelConfig(ctx, req.Keys); err != nil {
		return nil, err
	}
	return &admin.DelClientConfigResp{}, nil
}
