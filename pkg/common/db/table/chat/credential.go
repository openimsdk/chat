package chat

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
)

type Credential struct {
	UserID      string `bson:"user_id"`
	Account     string `bson:"account"`
	Type        int    `bson:"type"` // 1:phone;2:email
	AllowChange bool   `bson:"allow_change"`
}

func (Credential) TableName() string {
	return "credentials"
}

type CredentialInterface interface {
	Create(ctx context.Context, credential ...*Credential) error
	CreateOrUpdateAccount(ctx context.Context, credential *Credential) error
	Update(ctx context.Context, userID string, data map[string]any) error
	Find(ctx context.Context, userID string) ([]*Credential, error)
	FindAccount(ctx context.Context, accounts []string) ([]*Credential, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*Credential, error)
	TakeAccount(ctx context.Context, account string) (*Credential, error)
	Take(ctx context.Context, userID string) (*Credential, error)
	SearchNormalUser(ctx context.Context, keyword string, forbiddenID []string, pagination pagination.Pagination) (int64, []*Credential, error)
	SearchUser(ctx context.Context, keyword string, userIDs []string, pagination pagination.Pagination) (int64, []*Credential, error)
	Delete(ctx context.Context, userIDs []string) error
	DeleteByUserIDType(ctx context.Context, credentials ...*Credential) error
}
