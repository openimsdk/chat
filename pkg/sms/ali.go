package sms

import (
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	aliconf "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

func newAli() (SMS, error) {
	conf := &aliconf.Config{
		Endpoint:        tea.String(config.Config.VerifyCode.Ali.Endpoint),
		AccessKeyId:     tea.String(config.Config.VerifyCode.Ali.AccessKeyId),
		AccessKeySecret: tea.String(config.Config.VerifyCode.Ali.AccessKeySecret),
	}
	client, err := dysmsapi.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &ali{client: client}, nil
}

type ali struct {
	client *dysmsapi.Client
}

func (a *ali) Name() string {
	return "ali-sms"
}

func (a *ali) SendCode(ctx context.Context, areaCode string, phoneNumber string, verifyCode string) error {
	data, err := json.Marshal(&struct {
		Code string `json:"code"`
	}{Code: verifyCode})
	if err != nil {
		return errs.Wrap(err)
	}
	req := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(areaCode + phoneNumber),
		SignName:      tea.String(config.Config.VerifyCode.Ali.SignName),
		TemplateCode:  tea.String(config.Config.VerifyCode.Ali.VerificationCodeTemplateCode),
		TemplateParam: tea.String(string(data)),
	}
	_, err = a.client.SendSms(req)
	return errs.Wrap(err)
}
