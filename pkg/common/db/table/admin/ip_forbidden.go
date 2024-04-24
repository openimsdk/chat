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
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type IPForbidden struct {
	IP            string    `bson:"ip"`
	LimitRegister bool      `bson:"limit_register"`
	LimitLogin    bool      `bson:"limit_login"`
	CreateTime    time.Time `bson:"create_time"`
}

func (IPForbidden) IPForbidden() string {
	return "ip_forbiddens"
}

type IPForbiddenInterface interface {
	Take(ctx context.Context, ip string) (*IPForbidden, error)
	Find(ctx context.Context, ips []string) ([]*IPForbidden, error)
	Search(ctx context.Context, keyword string, state int32, pagination pagination.Pagination) (int64, []*IPForbidden, error)
	Create(ctx context.Context, ms []*IPForbidden) error
	Delete(ctx context.Context, ips []string) error
}
