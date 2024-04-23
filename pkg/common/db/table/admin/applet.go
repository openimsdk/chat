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

type Applet struct {
	ID         string    `bson:"id"`
	Name       string    `bson:"name"`
	AppID      string    `bson:"app_id"`
	Icon       string    `bson:"icon"`
	URL        string    `bson:"url"`
	MD5        string    `bson:"md5"`
	Size       int64     `bson:"size"`
	Version    string    `bson:"version"`
	Priority   uint32    `bson:"priority"`
	Status     uint8     `bson:"status"`
	CreateTime time.Time `bson:"create_time"`
}

func (Applet) TableName() string {
	return "applets"
}

type AppletInterface interface {
	Create(ctx context.Context, applets []*Applet) error
	Del(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, data map[string]any) error
	Take(ctx context.Context, id string) (*Applet, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*Applet, error)
	FindOnShelf(ctx context.Context) ([]*Applet, error)
	FindID(ctx context.Context, ids []string) ([]*Applet, error)
}
