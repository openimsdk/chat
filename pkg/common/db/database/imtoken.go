package database

import (
	"context"

	"github.com/openimsdk/chat/pkg/common/db/cache"
	"github.com/redis/go-redis/v9"
)

type IMTokenDatabase interface {
	GetIMToken(ctx context.Context, userID string) (string, error)
	SetIMToken(ctx context.Context, userID, token string) error
}

type imTokenDatabase struct {
	db cache.IMTokenInterface
}

func NewIMTokenDatabase(rdb redis.UniversalClient, expire int) IMTokenDatabase {
	return &imTokenDatabase{
		db: cache.NewIMTokenInterface(rdb, expire),
	}
}

func (i *imTokenDatabase) GetIMToken(ctx context.Context, userID string) (string, error) {
	return i.db.GetToken(ctx, userID)
}

func (i *imTokenDatabase) SetIMToken(ctx context.Context, userID, token string) error {
	return i.db.SetToken(ctx, userID, token)
}
