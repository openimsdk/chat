package organization

import (
	"context"
	"github.com/OpenIMSDK/chat/pkg/proto/organization"
	"time"
)

type OrganizationUser struct {
	UserID      string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname    string    `gorm:"column:nickname;size:256"`
	EnglishName string    `gorm:"column:english_name;size:256"`
	FaceURL     string    `gorm:"column:face_url;size:256"`
	Gender      int32     `gorm:"column:gender"` //1 ,2
	Station     string    `gorm:"column:station;size:256"`
	AreaCode    string    `gorm:"column:area_code;size:32"`
	Mobile      string    `gorm:"column:mobile;size:32"`
	Telephone   string    `gorm:"column:telephone;size:32"`
	Birth       time.Time `gorm:"column:birth"`
	Email       string    `gorm:"column:email;size:64"`
	Order       int32     `gorm:"column:order" json:"order"`
	Status      int32     `gorm:"column:status" json:"status"` //-1, 1
	CreateTime  time.Time `gorm:"column:create_time"`
	ChangeTime  time.Time `gorm:"column:change_time" json:"changeTime"`
}

type OrganizationUserInterface interface {
	Create(ctx context.Context, m *OrganizationUser) error
	Update(ctx context.Context, m *OrganizationUser) error
	Delete(ctx context.Context, userID string) error
	Get(ctx context.Context, userID string) (*OrganizationUser, error)
	SearchPage(ctx context.Context, positionList, userIDList []string, text string, sort []*organization.GetSearchUserListSort, pageNumber uint32, showNumber uint32) (uint32, []*OrganizationUser, error)
	GetNoDepartmentUserIDList(ctx context.Context) ([]string, error)
	GetList(ctx context.Context, userIDList []string) ([]*OrganizationUser, error)
	Search(ctx context.Context, positionList, userIDList []string, text string, sort []*organization.GetSearchUserListSort) ([]*OrganizationUser, error)
	GetPage(ctx context.Context, pageNumber, showNumber int) (int64, []*OrganizationUser, error)
}
