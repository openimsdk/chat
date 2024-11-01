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

package tokenverify

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/tools/errs"
)

const (
	TokenUser  = constant.NormalUser
	TokenAdmin = constant.AdminUser
)

type claims struct {
	UserID     string
	UserType   int32
	PlatformID int32
	jwt.RegisteredClaims
}

type Token struct {
	Expires time.Duration
	Secret  string
}

func (t *Token) secret() jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(t.Secret), nil
	}
}

func (t *Token) buildClaims(userID string, userType int32) claims {
	now := time.Now()
	return claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(t.Expires)),    // Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                   // Issuing time
			NotBefore: jwt.NewNumericDate(now.Add(-time.Minute)), // Begin Effective time
		},
	}
}

func (t *Token) getToken(str string) (string, int32, error) {
	token, err := jwt.ParseWithClaims(str, &claims{}, t.secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", 0, errs.ErrTokenMalformed.Wrap()
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", 0, errs.ErrTokenExpired.Wrap()
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", 0, errs.ErrTokenNotValidYet.Wrap()
			} else {
				return "", 0, errs.ErrTokenUnknown.Wrap()
			}
		} else {
			return "", 0, errs.ErrTokenNotValidYet.Wrap()
		}
	} else {
		claims, ok := token.Claims.(*claims)
		if claims.PlatformID != 0 {
			return "", 0, errs.ErrTokenExpired.Wrap()
		}
		if ok && token.Valid {
			return claims.UserID, claims.UserType, nil
		}
		return "", 0, errs.ErrTokenNotValidYet.Wrap()
	}
}

func (t *Token) CreateToken(UserID string, userType int32) (string, time.Duration, error) {
	if !(userType == TokenUser || userType == TokenAdmin) {
		return "", 0, errs.ErrTokenUnknown.WrapMsg("token type unknown")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.buildClaims(UserID, userType))
	str, err := token.SignedString([]byte(t.Secret))
	if err != nil {
		return "", 0, errs.Wrap(err)
	}
	return str, t.Expires, nil
}

func (t *Token) GetToken(token string) (string, int32, error) {
	userID, userType, err := t.getToken(token)
	if err != nil {
		return "", 0, err
	}
	if !(userType == TokenUser || userType == TokenAdmin) {
		return "", 0, errs.ErrTokenUnknown.WrapMsg("token type unknown")
	}
	return userID, userType, nil
}

//func (t *Token) GetAdminTokenCache(token string) (string, error) {
//	userID, userType, err := getToken(token)
//	if err != nil {
//		return "", err
//	}
//	if userType != TokenAdmin {
//		return "", errs.ErrTokenUnknown.WrapMsg("token type error")
//	}
//	return userID, nil
//}
//
//func (t *Token) GetUserToken(token string) (string, error) {
//	userID, userType, err := getToken(token)
//	if err != nil {
//		return "", err
//	}
//	if userType != TokenUser {
//		return "", errs.ErrTokenUnknown.WrapMsg("token type error")
//	}
//	return userID, nil
//}
