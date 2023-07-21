package apistruct

import "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"

type AdminLoginResp struct {
	AdminAccount string `json:"adminAccount"`
	AdminToken   string `json:"adminToken"`
	Nickname     string `json:"nickname"`
	FaceURL      string `json:"faceURL"`
	Level        int32  `json:"level"`
	AdminUserID  string `json:"adminUserID"`
	ImUserID     string
	ImToken      string
}

type SearchDefaultGroupResp struct {
	Total  uint32             `json:"total"`
	Groups []*GroupAttribute `json:"groups"`
}

type GroupAttribute struct {
	GroupID    string           `json:"groupID"`
	CreateTime int64            `json:"createTime"`
	Group      *sdkws.GroupInfo `json:"group"`
}
