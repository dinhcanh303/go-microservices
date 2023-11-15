package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/jordan-wright/email"
)

type emailSender struct {
	cfg *configs.Mail
}

// SendEmail implements EmailSender.
func (sender *emailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s<%s>", sender.cfg.FromName, sender.cfg.FromAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}
	smtpAuth := smtp.PlainAuth("", sender.cfg.Username, sender.cfg.Password, sender.cfg.Host)
	smtpServerAddress := fmt.Sprintf("%s:%s", sender.cfg.Host, sender.cfg.Port)
	if sender.cfg.Encryption == "tls" {
		return e.SendWithTLS(smtpServerAddress, smtpAuth, &tls.Config{
			InsecureSkipVerify: true,
		})
	}
	return e.Send(smtpServerAddress, smtpAuth)
}

var _ EmailSender = (*emailSender)(nil)

func NewEmailSender(cfg *configs.Mail) EmailSender {
	return &emailSender{
		cfg: cfg,
	}
}
