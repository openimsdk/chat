package botstruct

// see openim-sdk-core\sdk_struct\sdk_struct.go

type TextElem struct {
	Content string `json:"content"`
}

type AtElem struct {
	Text       string   `mapstructure:"text"`
	AtUserList []string `mapstructure:"atUserList" validate:"required,max=1000"`
	IsAtSelf   bool     `mapstructure:"isAtSelf"`
}
