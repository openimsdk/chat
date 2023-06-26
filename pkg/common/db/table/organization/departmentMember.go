package organization

import (
	"time"
)

type DepartmentMember struct {
	UserID          string     `gorm:"column:user_id;primary_key;size:64"`
	DepartmentID    string     `gorm:"column:department_id;primary_key;size:64"`
	Order           int32      `gorm:"column:order" json:"order"` //1,2
	Position        string     `gorm:"column:position;size:256" json:"position"`
	Leader          int32      `gorm:"column:leader" json:"leader"` //-1, 1
	Status          int32      `gorm:"column:status" json:"status"` //-1, 1
	EntryTime       time.Time  `gorm:"column:entry_time"`           // 入职时间
	TerminationTime *time.Time `gorm:"column:termination_time"`     // 离职时间
	CreateTime      time.Time  `gorm:"column:create_time"`
	ChangeTime      time.Time  `gorm:"column:change_time" json:"changeTime"`
}

type DepartmentMemberInterface interface {

}

