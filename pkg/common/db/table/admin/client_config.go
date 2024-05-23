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

import "context"

// ClientConfig config
type ClientConfig struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func (ClientConfig) TableName() string {
	return "client_config"
}

type ClientConfigInterface interface {
	Set(ctx context.Context, config map[string]string) error
	Get(ctx context.Context) (map[string]string, error)
	Del(ctx context.Context, keys []string) error
}
