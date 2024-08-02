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

package admin

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/chat/pkg/eerrs"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/log"
)

func (o *adminServer) CreateToken(ctx context.Context, req *adminpb.CreateTokenReq) (*adminpb.CreateTokenResp, error) {
	token, expire, err := o.Token.CreateToken(req.UserID, req.UserType)

	if err != nil {
		return nil, err
	}
	err = o.Database.CacheToken(ctx, req.UserID, token, expire)
	if err != nil {
		return nil, err
	}
	return &adminpb.CreateTokenResp{
		Token: token,
	}, nil
}

func (o *adminServer) ParseToken(ctx context.Context, req *adminpb.ParseTokenReq) (*adminpb.ParseTokenResp, error) {
	userID, userType, err := o.Token.GetToken(req.Token)
	if err != nil {
		return nil, err
	}
	m, err := o.Database.GetTokens(ctx, userID)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, eerrs.ErrTokenNotExist.Wrap()
	}
	if _, ok := m[req.Token]; !ok {
		return nil, eerrs.ErrTokenNotExist.Wrap()
	}

	return &adminpb.ParseTokenResp{
		UserID:   userID,
		UserType: userType,
	}, nil
}

func (o *adminServer) GetUserToken(ctx context.Context, req *adminpb.GetUserTokenReq) (*adminpb.GetUserTokenResp, error) {
	tokensMap, err := o.Database.GetTokens(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &adminpb.GetUserTokenResp{TokensMap: tokensMap}, nil
}

func (o *adminServer) InvalidateToken(ctx context.Context, req *adminpb.InvalidateTokenReq) (*adminpb.InvalidateTokenResp, error) {
	err := o.Database.DeleteToken(ctx, req.UserID)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	log.ZDebug(ctx, "delete token from redis", "userID", req.UserID)
	return &adminpb.InvalidateTokenResp{}, nil
}
