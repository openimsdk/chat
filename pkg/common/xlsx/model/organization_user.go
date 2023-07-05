package model

type OrganizationUser struct {
	UserID      string `column:"ID"`
	Nickname    string `column:"昵称"`
	EnglishName string `column:"英文名"`
	FaceURL     string `column:"头像"`
	Gender      string `column:"性别"`
	Station     string `column:"工位"`
	AreaCode    string `column:"区号"`
	Mobile      string `column:"手机号"`
	Telephone   string `column:"固定电话"`
	Birth       string `column:"生日"`
	Email       string `column:"邮箱"`
	Account     string `column:"账号"`
	Password    string `column:"密码"`
	Department  string `column:"所在部门"` //  开发/后端/Go:职位1;销售/后端/Go:职位2
}

func (OrganizationUser) SheetName() string {
	return "用户"
}
