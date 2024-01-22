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

package email

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestEmail(T *testing.T) {
	if err := InitConfig(); err != nil {
		panic(err)
	}
	tests := []struct {
		name string
		ctx  context.Context
		mail string
		code string
		want error
	}{
		{
			name: "success send email",
			ctx:  context.Background(),
			mail: "test@gmail.com",
			code: "5555",
			want: errors.New("nil"),
		},
		{
			name: "fail send email",
			ctx:  context.Background(),
			mail: "",
			code: "5555",
			want: errors.New("dial tcp :0: connectex: The requested address is not valid in its context."),
		},
	}
	mail := NewMail()

	for _, tt := range tests {
		T.Run(tt.name, func(t *testing.T) {
			if got := mail.SendMail(tt.ctx, tt.mail, tt.code); errors.Is(got, tt.want) {
				t.Errorf("%v have a err,%v", tt.name, tt.want)
			}
		})
	}
}

func InitConfig() error {
	yam, err := ioutil.ReadFile("../../config/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yam, &config.Config)
	if err != nil {
		return err
	}
	return nil
}
