// Copyright © 2023 OpenIM. All rights reserved.
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

package tx

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewMongo(client *mongo.Client) CtxTx {
	return &_Mongo{
		initialized: false,
		lock:        &sync.Mutex{},
		client:      client,
	}
}

type _Mongo struct {
	initialized bool
	lock        sync.Locker
	client      *mongo.Client
	tx          func(func(ctx context.Context) error) error
}

func (m *_Mongo) init(ctx context.Context) (err error) {
	m.lock.Lock()
	defer func() {
		if err == nil {
			m.initialized = true
		}
		m.lock.Unlock()
	}()
	if m.initialized {
		return nil
	}
	var res map[string]any
	if err := m.client.Database("admin").RunCommand(ctx, bson.M{"isMaster": 1}).Decode(&res); err != nil {
		return err
	}
	_, allowTx := res["setName"]
	if !allowTx {
		return nil
	}
	m.tx = func(fn func(ctx context.Context) error) error {
		sess, err := m.client.StartSession()
		if err != nil {
			return err
		}
		defer sess.EndSession(ctx)
		_, err = sess.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			return nil, fn(sessCtx)
		})
		return err
	}
	return nil
}

func (m *_Mongo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if !m.initialized {
		if err := m.init(ctx); err != nil {
			return err
		}
	}
	if m.tx == nil {
		return fn(ctx)
	}
	return m.tx(fn)
}
