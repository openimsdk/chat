package organization

import (
	"context"
	"time"
)

type Department struct {
	DepartmentID   string    `gorm:"column:department_id;primary_key;size:64" json:"departmentID"`
	FaceURL        string    `gorm:"column:face_url;size:255" json:"faceURL"`
	Name           string    `gorm:"column:name;size:256" json:"name" binding:"required"`
	ParentID       string    `gorm:"column:parent_id;size:64" json:"parentID" binding:"required"` // "0" or Real parent id
	Order          int32     `gorm:"column:order" json:"order" `                                  // 1, 2, ...
	DepartmentType int32     `gorm:"column:department_type" json:"departmentType"`                //1, 2...
	RelatedGroupID string    `gorm:"column:related_group_id;size:64" json:"relatedGroupID"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	ChangeTime     time.Time `gorm:"column:change_time" json:"changeTime"`
}

type DepartmentInterface interface {
	Create(ctx context.Context, department ...*Department) error
	FindOne(ctx context.Context, departmentID string) (*Department, error)
	Update(ctx context.Context, department *Department) error
	GetParent(ctx context.Context, id string) ([]*Department, error)
	GetList(ctx context.Context, departmentIdList []string) ([]*Department, error)
	UpdateParentID(ctx context.Context, oldParentID, newParentID string) error
	Delete(ctx context.Context, departmentIDList []string) error
	GetDepartment(ctx context.Context, departmentID string) (*Department, error)
	GetByName(ctx context.Context, name, parentID string) (*Department, error)
}
