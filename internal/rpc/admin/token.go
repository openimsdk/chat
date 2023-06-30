package admin

import (
	"context"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/tokenverify"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
)

func (*adminServer) CreateToken(ctx context.Context, req *admin.CreateTokenReq) (*admin.CreateTokenResp, error) {
	resp := &admin.CreateTokenResp{}
	var err error
	resp.Token, err = tokenverify.CreateToken(req.UserID, req.UserType, *config.Config.TokenPolicy.Expire)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (*adminServer) ParseToken(ctx context.Context, req *admin.ParseTokenReq) (*admin.ParseTokenResp, error) {
	resp := &admin.ParseTokenResp{}
	var err error
	resp.UserID, resp.UserType, err = tokenverify.GetToken(req.Token)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
