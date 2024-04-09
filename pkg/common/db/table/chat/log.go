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

package chat

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type Log struct {
	LogID      string    `bson:"log_id"`
	Platform   string    `bson:"platform"`
	UserID     string    `bson:"user_id"`
	CreateTime time.Time `bson:"create_time"`
	Url        string    `bson:"url"`
	FileName   string    `bson:"file_name"`
	SystemType string    `bson:"system_type"`
	Version    string    `bson:"version"`
	Ex         string    `bson:"ex"`
}

type LogInterface interface {
	Create(ctx context.Context, log []*Log) error
	Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*Log, error)
	Delete(ctx context.Context, logID []string, userID string) error
	Get(ctx context.Context, logIDs []string, userID string) ([]*Log, error)
}
