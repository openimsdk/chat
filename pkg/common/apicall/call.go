package apicall

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	constant2 "github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
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

func NewApiCaller[Req, Resp any](url string) ApiCaller[Req, Resp] {
	return Api[Req, Resp](url)
}

type Api[Req, Resp any] string

func (a Api[Req, Resp]) Call(ctx context.Context, req *Req) (*Resp, error) {
	log.ZInfo(ctx, "api req", "addr", string(a), "req", req)
	resp, err := a.call(ctx, req)
	if err != nil {
		log.ZError(ctx, "api resp", err)
		return nil, err
	}
	log.ZError(ctx, "api resp", err, "resp", resp)
	return resp, nil
}

func (a Api[Req, Resp]) call(ctx context.Context, req *Req) (*Resp, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, string(a), bytes.NewReader(reqBody))
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
	log.ZDebug(ctx, "call api successfully", "api", string(a))
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "read respBody successfully")
	var resp baseApiResponse[Resp]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, errs.NewCodeError(resp.ErrCode, resp.ErrMsg).WithDetail(resp.ErrDlt)
	}
	log.ZDebug(ctx, "unmarshal resp success")
	return resp.Data, nil
}
