package apistruct

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
