package chat

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
)

func ToDBAttributeUpdate(req *chat.UpdateUserInfoReq) (map[string]any, error) {
	update := make(map[string]any)
	if req.Account != nil {
		update["account"] = req.Account.Value
	}
	if req.AreaCode != nil {
		update["area_code"] = req.AreaCode.Value
	}
	if req.Email != nil {
		update["email"] = req.Email.Value
	}
	if req.Nickname != nil {
		if req.Nickname.Value == "" {
			return nil, errs.ErrArgs.Wrap("nickname can not be empty")
		}
		update["nickname"] = req.Nickname.Value
	}
	if req.FaceURL != nil {
		update["face_url"] = req.FaceURL.Value
	}
	if req.Gender != nil {
		update["gender"] = req.Gender.Value
	}
	if req.Level != nil {
		update["level"] = req.Level.Value
	}
	if req.Forbidden != nil {
		update["forbidden"] = req.Forbidden.Value
	}
	if req.Birth != nil {
		update["birth"] = req.Birth.Value
	}
	if req.AllowAddFriend != nil {
		update["allow_add_friend"] = req.AllowAddFriend.Value
	}
	if req.AllowBeep != nil {
		update["allow_beep"] = req.AllowBeep.Value
	}
	if req.AllowVibration != nil {
		update["allow_vibration"] = req.AllowVibration.Value
	}
	if len(update) == 0 {
		return nil, errs.ErrArgs.Wrap("no update info")
	}
	return update, nil
}
