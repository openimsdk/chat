package chat

import (
	"context"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/OpenIMSDK/tools/errs"
	"time"
)

func (o *chatSvr) UserLoginCount(ctx context.Context, req *chat.UserLoginCountReq) (*chat.UserLoginCountResp, error) {
	resp := &chat.UserLoginCountResp{}
	if req.Start > req.End {
		return nil, errs.ErrArgs.Wrap("start > end")
	}
	total, err := o.Database.NewUserCountTotal(ctx, nil)
	if err != nil {
		return nil, err
	}
	start := time.UnixMilli(req.Start)
	end := time.UnixMilli(req.End)
	count, loginCount, err := o.Database.UserLoginCountRangeEverydayTotal(ctx, &start, &end)
	if err != nil {
		return nil, err
	}
	resp.LoginCount = loginCount
	resp.UnloginCount = total - loginCount
	resp.Count = count
	return resp, nil
}
