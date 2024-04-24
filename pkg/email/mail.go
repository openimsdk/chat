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
	"github.com/openimsdk/tools/errs"
	"gopkg.in/gomail.v2"
)

type Mail interface {
	Name() string
	SendMail(ctx context.Context, mail string, verifyCode string) error
}

func NewMail(smtpAddr string, smtpPort int, senderMail, senderAuthorizationCode, title string) Mail {
	dail := gomail.NewDialer(smtpAddr, smtpPort, senderMail, senderAuthorizationCode)
	return &mail{
		title:      title,
		senderMail: senderMail,
		dail:       dail,
	}
}

type mail struct {
	senderMail string
	title      string
	dail       *gomail.Dialer
}

func (m *mail) Name() string {
	return "mail"
}

func (m *mail) SendMail(ctx context.Context, mail string, verifyCode string) error {
	msg := gomail.NewMessage()
	msg.SetHeader(`From`, m.senderMail)
	msg.SetHeader(`To`, []string{mail}...)
	msg.SetHeader(`Subject`, m.title)
	msg.SetBody(`text/html`, fmt.Sprintf("Your verification code is: %s. This code is valid for 5 minutes and should not be shared with others", verifyCode))
	return errs.Wrap(m.dail.DialAndSend(msg))
}
