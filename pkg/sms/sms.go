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

package sms

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"strings"
)

func New() (SMS, error) {
	switch strings.ToLower(config.Config.VerifyCode.Use) {
	case "":
		return empty{}, nil
	case "ali":
		return newAli()
	default:
		return nil, fmt.Errorf("not support sms: `%s`", config.Config.VerifyCode.Use)
	}
}

type SMS interface {
	Name() string
	SendCode(ctx context.Context, areaCode string, phoneNumber string, verifyCode string) error
}

type empty struct{}

func (e empty) Name() string {
	return "empty-sms"
}

func (e empty) SendCode(ctx context.Context, areaCode string, phoneNumber string, verifyCode string) error {
	return nil
}
