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
	"time"
)

// 禁止ip登录 注册.
type IPForbidden struct {
	IP            string    `gorm:"column:ip;primary_key;type:char(32)"`
	LimitRegister bool      `gorm:"column:limit_register"`
	LimitLogin    bool      `gorm:"column:limit_login"`
	CreateTime    time.Time `gorm:"column:create_time"`
}

func (IPForbidden) IPForbidden() string {
	return "ip_forbiddens"
}

type IPForbiddenInterface interface {
	NewTx(tx any) IPForbiddenInterface
	Take(ctx context.Context, ip string) (*IPForbidden, error)
	Find(ctx context.Context, ips []string) ([]*IPForbidden, error)
	Search(ctx context.Context, keyword string, state int32, page int32, size int32) (uint32, []*IPForbidden, error)
	Create(ctx context.Context, ms []*IPForbidden) error
	Delete(ctx context.Context, ips []string) error
}
