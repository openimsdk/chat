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

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/ormutil"
	"gorm.io/gorm"

	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
)

func NewIPForbidden(db *gorm.DB) admin.IPForbiddenInterface {
	return &IPForbidden{db: db}
}

type IPForbidden struct {
	db *gorm.DB
}

func (o *IPForbidden) NewTx(tx any) admin.IPForbiddenInterface {
	return &IPForbidden{db: tx.(*gorm.DB)}
}

func (o *IPForbidden) Take(ctx context.Context, ip string) (*admin.IPForbidden, error) {
	var f admin.IPForbidden
	return &f, errs.Wrap(o.db.WithContext(ctx).Where("ip = ?", ip).Take(&f).Error)
}

func (o *IPForbidden) Find(ctx context.Context, ips []string) ([]*admin.IPForbidden, error) {
	var forbiddens []*admin.IPForbidden
	return forbiddens, errs.Wrap(o.db.WithContext(ctx).Where("ip in ?", ips).Find(&forbiddens).Error)
}

func (o *IPForbidden) Search(ctx context.Context, keyword string, state int32, page int32, size int32) (uint32, []*admin.IPForbidden, error) {
	db := o.db.WithContext(ctx)
	switch state {
	case constant.LimitNil:
	case constant.LimitEmpty:
		db = db.Where("limit_register = 0 and limit_login = 0")
	case constant.LimitOnlyRegisterIP:
		db = db.Where("limit_register = 1 and limit_login = 0")
	case constant.LimitOnlyLoginIP:
		db = db.Where("limit_register = 0 and limit_login = 1")
	case constant.LimitRegisterIP:
		db = db.Where("limit_register = 1")
	case constant.LimitLoginIP:
		db = db.Where("limit_login = 1")
	case constant.LimitLoginRegisterIP:
		db = db.Where("limit_register = 1 and limit_login = 1")
	}
	return ormutil.GormSearch[admin.IPForbidden](db, []string{"ip"}, keyword, page, size)
}

func (o *IPForbidden) Create(ctx context.Context, ms []*admin.IPForbidden) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(&ms).Error)
}

func (o *IPForbidden) Delete(ctx context.Context, ips []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("ip in ?", ips).Delete(&admin.IPForbidden{}).Error)
}
