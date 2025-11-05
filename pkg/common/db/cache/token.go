package cache

import (
	"context"
	"time"

	"github.com/openimsdk/tools/utils/stringutil"

	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

const (
	chatToken = "CHAT_UID_TOKEN_STATUS:"
)

type TokenInterface interface {
	AddTokenFlag(ctx context.Context, userID string, token string, flag int) error
	AddTokenFlagNXEx(ctx context.Context, userID string, token string, flag int, expire time.Duration) (bool, error)
	GetTokensWithoutError(ctx context.Context, userID string) (map[string]int32, error)
	DeleteTokenByUid(ctx context.Context, userID string) error
}

type TokenCacheRedis struct {
	rdb          redis.UniversalClient
	accessExpire int64
}

func NewTokenInterface(rdb redis.UniversalClient) *TokenCacheRedis {
	return &TokenCacheRedis{rdb: rdb}
}

func (t *TokenCacheRedis) AddTokenFlag(ctx context.Context, userID string, token string, flag int) error {
	key := chatToken + userID
	return errs.Wrap(t.rdb.HSet(ctx, key, token, flag).Err())
}

func (t *TokenCacheRedis) AddTokenFlagNXEx(ctx context.Context, userID string, token string, flag int, expire time.Duration) (bool, error) {
	key := chatToken + userID
	isSet, err := t.rdb.HSetNX(ctx, key, token, flag).Result()
	if err != nil {
		return false, errs.Wrap(err)
	}
	if !isSet {
		// key already exists
		return false, nil
	}
	if err = t.rdb.Expire(ctx, key, expire).Err(); err != nil {
		return false, errs.Wrap(err)
	}
	return isSet, nil
}

func (t *TokenCacheRedis) GetTokensWithoutError(ctx context.Context, userID string) (map[string]int32, error) {
	key := chatToken + userID
	m, err := t.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int32)
	for k, v := range m {
		mm[k] = stringutil.StringToInt32(v)
	}
	return mm, nil
}

func (t *TokenCacheRedis) DeleteTokenByUid(ctx context.Context, userID string) error {
	key := chatToken + userID
	return errs.Wrap(t.rdb.Del(ctx, key).Err())
}
