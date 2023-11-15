package mail

type EmailSender interface {
	SendEmail(subject, content string, to, cc, bcc []string, attachFiles []string) error
}
