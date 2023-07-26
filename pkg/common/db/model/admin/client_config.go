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
	"gorm.io/gorm"

	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
)

func NewClientConfig(db *gorm.DB) admin.ClientConfigInterface {
	return &ClientConfig{db: db}
}

type ClientConfig struct {
	db *gorm.DB
}

func (o *ClientConfig) NewTx(tx any) admin.ClientConfigInterface {
	return &ClientConfig{db: tx.(*gorm.DB)}
}

func (o *ClientConfig) Set(ctx context.Context, config map[string]string) error {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range config {
			var cc admin.ClientConfig
			if err := tx.Where("`key` = ?", key).Take(&cc).Error; err == nil {
				if cc.Value == value {
					continue
				}
				if err := tx.Where("`key` = ?", key).Model(&admin.ClientConfig{}).Update("value", value).Error; err != nil {
					return err
				}
			} else if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&admin.ClientConfig{Key: key, Value: value}).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})
	return errs.Wrap(err)
}

func (o *ClientConfig) Del(ctx context.Context, keys []string) error {
	return errs.Wrap(o.db.WithContext(ctx).Where("`key` in ?", keys).Delete(&admin.ClientConfig{}).Error)
}

func (o *ClientConfig) Get(ctx context.Context) (map[string]string, error) {
	var cs []*admin.ClientConfig
	if err := o.db.WithContext(ctx).Find(&cs).Error; err != nil {
		return nil, err
	}
	cm := make(map[string]string)
	for _, config := range cs {
		cm[config.Key] = config.Value
	}
	return cm, nil
}
