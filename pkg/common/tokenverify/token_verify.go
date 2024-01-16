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

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/golang-jwt/jwt/v4"
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

func buildClaims(userID string, userType int32, ttl int64) claims {
	now := time.Now()
	before := now.Add(-time.Minute * 5)
	return claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*24) * time.Hour)), // Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                        // Issuing time
			NotBefore: jwt.NewNumericDate(before),                                     // Begin Effective time
		},
	}
}

func CreateToken(UserID string, userType int32, ttl int64) (string, error) {
	if !(userType == TokenUser || userType == TokenAdmin) {
		return "", errs.ErrTokenUnknown.Wrap("token type unknown")
	}
	claims := buildClaims(UserID, userType, ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(*config.Config.Secret))
	if err != nil {
		return "", errs.Wrap(err, "")
	}
	return tokenString, nil
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(*config.Config.Secret), nil
	}
}

func getToken(t string) (string, int32, error) {
	token, err := jwt.ParseWithClaims(t, &claims{}, secret())
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
			return "", 0, errs.ErrTokenNotExist.Wrap()
		}
		if ok && token.Valid {
			return claims.UserID, claims.UserType, nil
		}
		return "", 0, errs.ErrTokenNotValidYet.Wrap()
	}
}

func GetToken(token string) (string, int32, error) {
	userID, userType, err := getToken(token)
	if err != nil {
		return "", 0, err
	}
	if !(userType == TokenUser || userType == TokenAdmin) {
		return "", 0, errs.ErrTokenUnknown.Wrap("token type unknown")
	}
	return userID, userType, nil
}

func GetAdminToken(token string) (string, error) {
	userID, userType, err := getToken(token)
	if err != nil {
		return "", err
	}
	if userType != TokenAdmin {
		return "", errs.ErrTokenInvalid.Wrap("token type error")
	}
	return userID, nil
}

func GetUserToken(token string) (string, error) {
	userID, userType, err := getToken(token)
	if err != nil {
		return "", err
	}
	if userType != TokenUser {
		return "", errs.ErrTokenInvalid.Wrap("token type error")
	}
	return userID, nil
}
