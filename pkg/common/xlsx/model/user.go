package model

type User struct {
	UserID      string `column:"user_id"`
	Nickname    string `column:"nickname"`
	FaceURL     string `column:"face_url"`
	Birth       string `column:"birth"`
	Gender      string `column:"gender"`
	AreaCode    string `column:"area_code"`
	PhoneNumber string `column:"phone_number"`
	Email       string `column:"email"`
	Account     string `column:"account"`
	Password    string `column:"password"`
}

func (User) SheetName() string {
	return "user"
}
