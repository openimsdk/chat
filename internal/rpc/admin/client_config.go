package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
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
		return nil, errs.ErrArgs.Wrap("update config empty")
	}
	conf := make(map[string]*string)
	for key, value := range req.Config {
		if value == nil {
			conf[key] = nil
		} else {
			temp := value.Value
			conf[key] = &temp
		}
	}
	if err := o.Database.SetConfig(ctx, conf); err != nil {
		return nil, err
	}
	return &admin.SetClientConfigResp{}, nil
}
