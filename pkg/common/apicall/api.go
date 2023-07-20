package apicall

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	constant2 "github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/chat/pkg/common/config"
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
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, config.Config.OpenIM_url+string(a), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	operationID := utils.ToString(ctx.Value(constant2.OperationID))
	request.Header.Set(constant2.OperationID, operationID)
	if token := ctx.Value(constant2.Token); token != nil {
		request.Header.Set(constant2.Token, utils.ToString(ctx.Value(constant2.Token)))
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

var apiURL = config.Config.OpenIM_url

var (
	UserToken = NewApiCaller[auth.UserTokenReq, auth.UserTokenResp](apiURL + "/auth/user_token")
)

func test() error {
	ctx := context.Background()
	resp, err := UserToken.Call(ctx, &auth.UserTokenReq{
		Secret:     "",
		PlatformID: 0,
		UserID:     "",
	})
	if err != nil {
		return err
	}
	_ = resp.Token
	_ = resp.ExpireTimeSeconds

	return nil
}
