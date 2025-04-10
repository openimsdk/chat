// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
