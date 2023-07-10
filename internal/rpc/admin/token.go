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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/tokenverify"
	"github.com/OpenIMSDK/chat/pkg/proto/admin"
)

func (*adminServer) CreateToken(ctx context.Context, req *admin.CreateTokenReq) (*admin.CreateTokenResp, error) {
	defer log.ZDebug(ctx, "return")
	token, err := tokenverify.CreateToken(req.UserID, req.UserType, *config.Config.TokenPolicy.Expire)
	if err != nil {
		return nil, err
	}
	return &admin.CreateTokenResp{
		Token: token,
	}, nil
}

func (*adminServer) ParseToken(ctx context.Context, req *admin.ParseTokenReq) (*admin.ParseTokenResp, error) {
	userID, userType, err := tokenverify.GetToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &admin.ParseTokenResp{
		UserID:   userID,
		UserType: userType,
	}, nil
}
