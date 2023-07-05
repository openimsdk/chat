package model

type Department struct {
	DepartmentID   string `column:"ID"`
	Name           string `column:"名字"`
	FaceURL        string `column:"头像"`
	Parent         string `column:"上级部门"` //    开发/后端/Go
	RelatedGroupID string `column:"相关组"`
	DepartmentType int32  `column:"类型"`
}

func (Department) SheetName() string {
	return "部门"
}
