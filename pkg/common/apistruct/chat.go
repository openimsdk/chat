package apistruct

type UserRegisterResp struct {
	ImToken   string
	ChatToken string
	UserID    string
}

type LoginResp struct {
	ImToken   string
	ChatToken string
	UserID    string
}

type UpdateUserInfoResp struct{}
