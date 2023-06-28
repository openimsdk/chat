package chat

import (
	"context"
	"time"
)

// Attribute 用户属性表
type Attribute struct {
	UserID         string    `gorm:"column:user_id;primary_key;type:char(64)"`
	Account        string    `gorm:"column:account;type:char(64)"`
	PhoneNumber    string    `gorm:"column:phone_number;type:varchar(32)"`
	AreaCode       string    `gorm:"column:area_code;type:varchar(8)"`
	Email          string    `gorm:"column:email;type:varchar(64)" `
	Nickname       string    `gorm:"column:nickname;type:varchar(64)" `
	FaceURL        string    `gorm:"column:face_url;type:varchar(255)" `
	Gender         int32     `gorm:"column:gender"`
	CreateTime     time.Time `gorm:"column:create_time"`
	ChangeTime     time.Time `gorm:"column:change_time"`
	BirthTime      time.Time `gorm:"column:birth_time"`
	Level          int32     `gorm:"column:level;default:1"`
	AllowVibration int32     `gorm:"column:allow_vibration;default:1"`
	AllowBeep      int32     `gorm:"column:allow_beep;default:1"`
	AllowAddFriend int32     `gorm:"column:allow_add_friend;default:1"`
}

func (Attribute) TableName() string {
	return "attributes"
}

type AttributeInterface interface {
	NewTx(tx any) AttributeInterface
	Create(ctx context.Context, attribute ...*Attribute) error
	Update(ctx context.Context, userID string, data map[string]any) error
	Find(ctx context.Context, userIds []string) ([]*Attribute, error)
	FindAccount(ctx context.Context, accounts []string) ([]*Attribute, error)
	Search(ctx context.Context, keyword string, genders []int32, page int32, size int32) (uint32, []*Attribute, error)
	TakePhone(ctx context.Context, areaCode string, phoneNumber string) (*Attribute, error)
	TakeAccount(ctx context.Context, account string) (*Attribute, error)
	Take(ctx context.Context, userID string) (*Attribute, error)
	GetAccountList(ctx context.Context, accountList []string) ([]*Attribute, error)
	ExistPhoneNumber(ctx context.Context, areaCode, phoneNumber string) (bool, error)
}
