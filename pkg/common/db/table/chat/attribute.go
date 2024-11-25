package chat

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/pagination"
)

type Attribute struct {
	UserID           string    `bson:"user_id"`
	Account          string    `bson:"account"`
	PhoneNumber      string    `bson:"phone_number"`
	AreaCode         string    `bson:"area_code"`
	Email            string    `bson:"email"`
	Nickname         string    `bson:"nickname"`
	FaceURL          string    `bson:"face_url"`
	Gender           int32     `bson:"gender"`
	CreateTime       time.Time `bson:"create_time"`
	ChangeTime       time.Time `bson:"change_time"`
	BirthTime        time.Time `bson:"birth_time"`
	Level            int32     `bson:"level"`
	AllowVibration   int32     `bson:"allow_vibration"`
	AllowBeep        int32     `bson:"allow_beep"`
	AllowAddFriend   int32     `bson:"allow_add_friend"`
	GlobalRecvMsgOpt int32     `bson:"global_recv_msg_opt"`
	RegisterType     int32     `bson:"register_type"`
}

func (Attribute) TableName() string {
	return "attributes"
}

type AttributeInterface interface {
	// NewTx(tx any) AttributeInterface
	Create(ctx context.Context, attribute ...*Attribute) error
	Update(ctx context.Context, userID string, data map[string]any) error
	Find(ctx context.Context, userIds []string) ([]*Attribute, error)
	FindAccount(ctx context.Context, accounts []string) ([]*Attribute, error)
	Search(ctx context.Context, keyword string, genders []int32, pagination pagination.Pagination) (int64, []*Attribute, error)
	TakePhone(ctx context.Context, areaCode string, phoneNumber string) (*Attribute, error)
	TakeEmail(ctx context.Context, email string) (*Attribute, error)
	TakeAccount(ctx context.Context, account string) (*Attribute, error)
	Take(ctx context.Context, userID string) (*Attribute, error)
	SearchNormalUser(ctx context.Context, keyword string, forbiddenID []string, gender int32, pagination pagination.Pagination) (int64, []*Attribute, error)
	SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pagination pagination.Pagination) (int64, []*Attribute, error)
	Delete(ctx context.Context, userIDs []string) error
}
