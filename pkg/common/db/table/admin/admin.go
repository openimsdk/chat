package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

// Admin user
type Admin struct {
	Account    string    `bson:"account"`
	Password   string    `bson:"password"`
	FaceURL    string    `bson:"face_url"`
	Nickname   string    `bson:"nickname"`
	UserID     string    `bson:"user_id"`
	Level      int32     `bson:"level"`
	CreateTime time.Time `bson:"create_time"`
}

func (Admin) TableName() string {
	return "admins"
}

type AdminInterface interface {
	Create(ctx context.Context, admins []*Admin) error
	Take(ctx context.Context, account string) (*Admin, error)
	TakeUserID(ctx context.Context, userID string) (*Admin, error)
	Update(ctx context.Context, account string, update map[string]any) error
	ChangePassword(ctx context.Context, userID string, newPassword string) error
	Delete(ctx context.Context, userIDs []string) error
	Search(ctx context.Context, pagination pagination.Pagination) (int64, []*Admin, error)
}
