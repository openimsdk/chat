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

package imapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/openimsdk/chat/pkg/common/constant"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"gorm.io/gorm/utils"
)

type baseApiResponse[T any] struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    *T     `json:"data"`
}

var client = &http.Client{
	Timeout: time.Second * 10,
}

type ApiCaller[Req, Resp any] interface {
	Call(ctx context.Context, apiPrefix string, req *Req) (*Resp, error)
}

func NewApiCaller[Req, Resp any](api string) ApiCaller[Req, Resp] {
	return &caller[Req, Resp]{
		api: api,
	}
}

type caller[Req, Resp any] struct {
	api string
}

func (a caller[Req, Resp]) Call(ctx context.Context, apiPrefix string, req *Req) (*Resp, error) {
	resp, err := a.call(ctx, apiPrefix, req)
	if err != nil {
		log.ZError(ctx, "caller resp", err)
		return nil, err
	}
	log.ZInfo(ctx, "resp", resp)
	return resp, nil
}

func (a caller[Req, Resp]) call(ctx context.Context, apiPrefix string, req *Req) (*Resp, error) {
	url := apiPrefix + a.api
	defer func(start time.Time) {
		log.ZDebug(ctx, "api call caller time", "api", a.api, "cost", time.Since(start).String())
	}(time.Now())
	log.ZInfo(ctx, "caller req", "addr", url, "req", req)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	operationID := utils.ToString(ctx.Value(pconstant.OperationID))
	request.Header.Set(pconstant.OperationID, operationID)
	if token, _ := ctx.Value(constant.CtxApiToken).(string); token != "" {
		request.Header.Set(pconstant.Token, token)
		log.ZDebug(ctx, "req token", "token", token)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "call caller successfully", "code", response.Status)
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errs.WrapMsg(err, "read http response body", "url", url, "code", response.StatusCode)
	}
	log.ZDebug(ctx, "read respBody successfully", "code", "url", url, response.StatusCode, "body", string(data))
	var resp baseApiResponse[Resp]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, errs.WrapMsg(err, string(data))
	}
	if resp.ErrCode != 0 {
		return nil, errs.NewCodeError(resp.ErrCode, resp.ErrMsg).WithDetail(resp.ErrDlt).Wrap()
	}
	return resp.Data, nil
}
