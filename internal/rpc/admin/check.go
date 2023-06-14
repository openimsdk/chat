package admin

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/chat/pkg/common/db/dbutil"
	"github.com/OpenIMSDK/chat/pkg/eerrs"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
)

func (o *adminServer) CheckRegisterForbidden(ctx context.Context, req *admin.CheckRegisterForbiddenReq) (*admin.CheckRegisterForbiddenResp, error) {
	forbiddens, err := o.Database.FindIPForbidden(ctx, []string{req.Ip})
	if err != nil {
		return nil, err
	}
	for _, forbidden := range forbiddens {
		if forbidden.LimitRegister {
			return nil, eerrs.ErrForbidden.Wrap()
		}
	}
	return &admin.CheckRegisterForbiddenResp{}, nil
}

func (o *adminServer) CheckLoginForbidden(ctx context.Context, req *admin.CheckLoginForbiddenReq) (*admin.CheckLoginForbiddenResp, error) {
	forbiddens, err := o.Database.FindIPForbidden(ctx, []string{req.Ip})
	if err != nil {
		return nil, err
	}
	for _, forbidden := range forbiddens {
		if forbidden.LimitLogin {
			return nil, eerrs.ErrForbidden.Wrap("ip forbidden")
		}
	}
	if _, err := o.Database.GetLimitUserLoginIP(ctx, req.UserID, req.Ip); err != nil {
		if !dbutil.IsGormNotFound(err) {
			return nil, err
		}
		count, err := o.Database.CountLimitUserLoginIP(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, eerrs.ErrForbidden.Wrap("user ip forbidden")
		}
	}
	if forbiddenAccount, err := o.Database.GetBlockInfo(ctx, req.UserID); err == nil {
		return nil, eerrs.ErrForbidden.Wrap(fmt.Sprintf("account forbidden: %s", forbiddenAccount.Reason))
	} else if !dbutil.IsGormNotFound(err) {
		return nil, err
	}
	return &admin.CheckLoginForbiddenResp{}, nil
}
