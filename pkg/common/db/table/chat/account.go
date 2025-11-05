package chat

import (
	"context"
	"time"
)

type Account struct {
	UserID         string    `bson:"user_id"`
	Password       string    `bson:"password"`
	CreateTime     time.Time `bson:"create_time"`
	ChangeTime     time.Time `bson:"change_time"`
	OperatorUserID string    `bson:"operator_user_id"`
}

func (Account) TableName() string {
	return "accounts"
}

type AccountInterface interface {
	Create(ctx context.Context, accounts ...*Account) error
	Take(ctx context.Context, userId string) (*Account, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	UpdatePassword(ctx context.Context, userId string, password string) error
	Delete(ctx context.Context, userIDs []string) error
}
