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
	emailFromName := sender.cfg.FromName
	emailFromAddress := sender.cfg.FromAddress
	emailPassword := sender.cfg.Password
	emailHost := sender.cfg.Host
	emailPort := sender.cfg.Port
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s<%s>", emailFromName, emailFromAddress)
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
	smtpAuth := smtp.PlainAuth("", emailFromAddress, emailPassword, emailHost)
	smtpServerAddress := fmt.Sprintf("%s:%s", emailHost, emailPort)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         emailHost,
	}
	if sender.cfg.Encryption == "tls" {
		return e.SendWithStartTLS(smtpServerAddress, smtpAuth, tlsConfig)
	}
	return e.Send(smtpServerAddress, smtpAuth)
}

var _ EmailSender = (*emailSender)(nil)

func NewEmailSender(cfg *configs.Mail) EmailSender {
	return &emailSender{
		cfg: cfg,
	}
}
