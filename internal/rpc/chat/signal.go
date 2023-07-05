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
	"errors"

	"github.com/OpenIMSDK/chat/pkg/proto/chat"
)

func (*chatSvr) AddSignalRecord(ctx context.Context, req *chat.AddSignalRecordReq) (*chat.AddSignalRecordResp, error) {
	return nil, errors.New("todo AddSignalRecord")
}

func (*chatSvr) GetSignalRecords(ctx context.Context, req *chat.GetSignalRecordsReq) (*chat.GetSignalRecordsResp, error) {
	return nil, errors.New("todo GetSignalRecords")
}
