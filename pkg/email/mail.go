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
	"fmt"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/tools/errs"
	"gopkg.in/gomail.v2"
)

func NewMail() Mail {
	dail := gomail.NewDialer(
		config.Config.VerifyCode.Mail.SmtpAddr,
		config.Config.VerifyCode.Mail.SmtpPort,
		config.Config.VerifyCode.Mail.SenderMail,
		config.Config.VerifyCode.Mail.SenderAuthorizationCode)

	return &mail{dail: dail}
}

type Mail interface {
	Name() string
	SendMail(ctx context.Context, mail string, verifyCode string) error
}

type mail struct {
	dail *gomail.Dialer
}

func (a *mail) Name() string {
	return "mail"
}

func (a *mail) SendMail(ctx context.Context, mail string, verifyCode string) error {
	m := gomail.NewMessage()
	m.SetHeader(`From`, config.Config.VerifyCode.Mail.SenderMail)
	m.SetHeader(`To`, []string{mail}...)
	m.SetHeader(`Subject`, config.Config.VerifyCode.Mail.Title)
	m.SetBody(`text/html`, fmt.Sprintf("您的验证码为:%s，该验证码5分钟内有效，请勿泄露于他人。", verifyCode))

	// Send
	err := a.dail.DialAndSend(m)
	return errs.Wrap(err)
}
