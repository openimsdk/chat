package cache

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/redis/go-redis/v9"
)

const (
	chatToken = "CHAT_UID_TOKEN_STATUS:"
)

type TokenInterface interface {
	AddTokenFlag(ctx context.Context, userID string, token string, flag int) error
	GetTokensWithoutError(ctx context.Context, userID string) (map[string]int32, error)
}

type TokenCacheRedis struct {
	rdb redis.UniversalClient
}

func NewTokenInterface(rdb redis.UniversalClient) *TokenCacheRedis {
	return &TokenCacheRedis{rdb: rdb}
}

func (t *TokenCacheRedis) AddTokenFlag(ctx context.Context, userID string, token string, flag int) error {
	key := chatToken + userID
	return errs.Wrap(t.rdb.HSet(ctx, key, token, flag).Err())
}

func (t *TokenCacheRedis) GetTokensWithoutError(ctx context.Context, userID string) (map[string]int32, error) {
	key := chatToken + userID
	m, err := t.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int32)
	for k, v := range m {
		mm[k] = utils.StringToInt32(v)
	}
	return mm, nil
}
