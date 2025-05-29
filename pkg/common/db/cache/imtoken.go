package cache

import (
	"context"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

const (
	chatPrefix = "CHAT:"
	imToken    = chatPrefix + "IM_TOKEN:"
)

func getIMTokenKey(userID string) string {
	return imToken + userID
}

type IMTokenInterface interface {
	GetToken(ctx context.Context, userID string) (string, error)
	SetToken(ctx context.Context, userID, token string) error
}

type imTokenCacheRedis struct {
	rdb    redis.UniversalClient
	expire time.Duration
}

func NewIMTokenInterface(rdb redis.UniversalClient, expire int) IMTokenInterface {
	return &imTokenCacheRedis{rdb: rdb, expire: time.Duration(expire) * time.Minute}
}

func (i *imTokenCacheRedis) GetToken(ctx context.Context, userID string) (string, error) {
	key := getIMTokenKey(userID)
	token, err := i.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return token, nil
}

func (i *imTokenCacheRedis) SetToken(ctx context.Context, userID, token string) error {
	key := getIMTokenKey(userID)
	err := i.rdb.Set(ctx, key, token, i.expire).Err()
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}
