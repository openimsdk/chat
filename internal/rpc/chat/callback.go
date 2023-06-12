package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/callbackstruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	constant2 "github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/eerrs"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
)

func (o *chatSvr) OpenIMCallback(ctx context.Context, req *chat.OpenIMCallbackReq) (*chat.OpenIMCallbackResp, error) {
	switch req.Command {
	case constant.CallbackBeforeAddFriendCommand:
		var data callbackstruct.CallbackBeforeAddFriendReq
		if err := json.Unmarshal([]byte(req.Body), &data); err != nil {
			return nil, errs.Wrap(err)
		}
		user, err := o.Database.GetAttribute(ctx, data.ToUserID)
		if err != nil {
			return nil, err
		}
		log.ZInfo(ctx, "OpenIMCallback", "user", user)
		if user.AllowAddFriend != constant2.OrdinaryUserAddFriendEnable {
			return nil, eerrs.ErrRefuseFriend.Wrap(fmt.Sprintf("state %d", user.AllowAddFriend))
		}
		return &chat.OpenIMCallbackResp{}, nil
	default:
		return nil, errs.ErrArgs.Wrap(fmt.Sprintf("invalid command %s", req.Command))
	}
}
