package apicall

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	constant2 "github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"gorm.io/gorm/utils"
	"io"
	"net/http"
)

type baseApiResponse[T any] struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    *T     `json:"data"`
}

type ApiCaller[Req, Resp any] interface {
	Call(ctx context.Context, req *Req) (*Resp, error)
}

func NewApiCaller[Req, Resp any](api string, prefix func() string) ApiCaller[Req, Resp] {
	return &caller[Req, Resp]{
		api:    api,
		prefix: prefix,
	}
}

type caller[Req, Resp any] struct {
	api    string
	prefix func() string
}

func (a caller[Req, Resp]) Call(ctx context.Context, req *Req) (*Resp, error) {
	resp, err := a.call(ctx, req)
	if err != nil {
		log.ZError(ctx, "caller resp", err)
		return nil, err
	}
	log.ZInfo(ctx, "caller resp", err, "resp", resp)
	return resp, nil
}

func (a caller[Req, Resp]) call(ctx context.Context, req *Req) (*Resp, error) {
	url := a.prefix() + a.api
	log.ZInfo(ctx, "caller req", "addr", url, "req", req)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	operationID := utils.ToString(ctx.Value(constant2.OperationID))
	request.Header.Set(constant2.OperationID, operationID)
	if token, _ := ctx.Value(constant.CtxApiToken).(string); token != "" {
		request.Header.Set(constant2.Token, token)
		log.ZDebug(ctx, "req token", "token", token)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "call caller successfully", "code", response.Status)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "read respBody successfully", "body", string(data))
	var resp baseApiResponse[Resp]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, errs.NewCodeError(resp.ErrCode, resp.ErrMsg).WithDetail(resp.ErrDlt)
	}
	return resp.Data, nil
}
