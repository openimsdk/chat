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
