package bot

import "github.com/openimsdk/chat/pkg/protocol/bot"

func ToDBAgentUpdate(req *bot.UpdateAgentReq) map[string]any {
	update := make(map[string]any)
	if req.Key != nil {
		update["key"] = req.Key
	}
	if req.FaceURL != nil {
		update["face_url"] = req.FaceURL
	}
	if req.NickName != nil {
		update["nick_name"] = req.NickName
	}
	if req.Identity != nil {
		update["identity"] = req.Identity
	}
	if req.Url != nil {
		update["url"] = req.Url
	}
	return update
}

func ToDBConversationRespIDUpdate(respID string) map[string]any {
	update := map[string]any{
		"previous_response_id": respID,
	}
	return update
}
