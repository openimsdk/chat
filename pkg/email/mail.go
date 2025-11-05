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
