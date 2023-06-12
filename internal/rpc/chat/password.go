package chat

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/mctx"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
)

func (o *chatSvr) ResetPassword(ctx context.Context, req *chat.ResetPasswordReq) (*chat.ResetPasswordResp, error) {
	if req.Password == "" {
		return nil, errs.ErrArgs.Wrap("password must be set")
	}
	verifyCodeID, err := o.verifyCode(ctx, o.verifyCodeJoin(req.VerifyCode, req.PhoneNumber), req.VerifyCode)
	if err != nil {
		return nil, err
	}
	attribute, err := o.Database.GetAttributeByPhone(ctx, req.AreaCode, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	err = o.Database.UpdatePasswordAndDeleteVerifyCode(ctx, attribute.UserID, req.Password, verifyCodeID)
	if err != nil {
		return nil, err
	}
	return &chat.ResetPasswordResp{}, nil
}

func (o *chatSvr) ChangePassword(ctx context.Context, req *chat.ChangePasswordReq) (*chat.ChangePasswordResp, error) {
	if req.Password == "" {
		return nil, errs.ErrArgs.Wrap("new password must be set")
	}
	opUserID, userType, err := mctx.Check(ctx)
	if err != nil {
		return nil, err
	}
	switch userType {
	case constant.NormalUser:
		if req.UserID == "" {
			req.UserID = opUserID
		}
		if req.UserID != opUserID {
			return nil, errs.ErrNoPermission.Wrap("no permission change other user password")
		}
	case constant.AdminUser:
		if req.UserID == "" {
			return nil, errs.ErrArgs.Wrap("user id must be set")
		}
	default:
		return nil, errs.ErrInternalServer.Wrap("invalid user type")
	}
	_, err = o.Database.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if err := o.Database.UpdatePassword(ctx, req.UserID, req.Password); err != nil {
		return nil, err
	}
	return &chat.ChangePasswordResp{}, nil
}
