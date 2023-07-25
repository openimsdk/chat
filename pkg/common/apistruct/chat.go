package apistruct

type UserRegisterResp struct {
	ImToken   string `json:"imToken"`
	ChatToken string `json:"chatToken"`
	UserID    string `json:"userID"`
}

type LoginResp struct {
	ImToken   string `json:"imToken"`
	ChatToken string `json:"chatToken"`
	UserID    string `json:"userID"`
}

type UpdateUserInfoResp struct{}
