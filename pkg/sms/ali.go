package sms

import (
	"context"
	"encoding/json"

	aliconf "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/openimsdk/tools/errs"
)

func NewAli(endpoint, accessKeyId, accessKeySecret, signName, verificationCodeTemplateCode string) (SMS, error) {
	conf := &aliconf.Config{
		Endpoint:        tea.String(endpoint),
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	client, err := dysmsapi.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &ali{
		signName:                     signName,
		verificationCodeTemplateCode: verificationCodeTemplateCode,
		client:                       client,
	}, nil
}

type ali struct {
	signName                     string
	verificationCodeTemplateCode string
	client                       *dysmsapi.Client
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
		SignName:      tea.String(a.signName),
		TemplateCode:  tea.String(a.verificationCodeTemplateCode),
		TemplateParam: tea.String(string(data)),
	}
	_, err = a.client.SendSms(req)
	return errs.Wrap(err)
}
