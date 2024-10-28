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
	"github.com/openimsdk/chat/pkg/common/tokenverify"
	"github.com/openimsdk/tools/utils/stringutil"
	"sort"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

const (
	chatToken       = "CHAT_UID_TOKEN_STATUS:"
	userMaxTokenNum = 10
)

type TokenInterface interface {
	SetTokenExpire(ctx context.Context, userID string, token string, expire time.Duration) error
	GetTokensWithoutError(ctx context.Context, userID string) (map[string]int32, error)
	DeleteTokenByUid(ctx context.Context, userID string) error
}

type TokenCacheRedis struct {
	token        *tokenverify.Token
	rdb          redis.UniversalClient
	accessExpire int64
}

func NewTokenInterface(rdb redis.UniversalClient, token *tokenverify.Token) *TokenCacheRedis {
	return &TokenCacheRedis{rdb: rdb, token: token}
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

func (t *TokenCacheRedis) SetTokenExpire(ctx context.Context, userID string, token string, expire time.Duration) error {
	key := chatToken + userID
	if err := t.rdb.HSet(ctx, key, token, "0").Err(); err != nil {
		return errs.Wrap(err)
	}
	if err := t.rdb.Expire(ctx, key, expire).Err(); err != nil {
		return errs.Wrap(err)
	}
	mm, err := t.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return errs.Wrap(err)
	}
	if len(mm) <= 1 {
		return nil
	}
	var (
		fields []string
		ts     tokenTimes
	)
	now := time.Now()
	for k := range mm {
		if k == token {
			continue
		}
		val := t.token.GetExpire(k)
		if val.IsZero() || val.Before(now) {
			fields = append(fields, k)
		} else {
			ts = append(ts, tokenTime{Token: k, Time: val})
		}
	}
	var sorted bool
	var index int
	for i := len(mm) - len(fields); i > userMaxTokenNum; i-- {
		if !sorted {
			sorted = true
			sort.Sort(ts)
		}
		fields = append(fields, ts[index].Token)
		index++
	}
	if len(fields) > 0 {
		if err := t.rdb.HDel(ctx, key, fields...).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

type tokenTime struct {
	Token string
	Time  time.Time
}

type tokenTimes []tokenTime

func (t tokenTimes) Len() int {
	return len(t)
}

func (t tokenTimes) Less(i, j int) bool {
	return t[i].Time.Before(t[j].Time)
}

func (t tokenTimes) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
