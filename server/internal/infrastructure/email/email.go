package email

import (
	"server/internal/config"
	"server/pkg/logger"

	"github.com/mailgun/mailgun-go"
)

type IEmailService interface {
	SendEmail(email, subject, body, to string) error
}

type emailServiceImpl struct {
	client *mailgun.MailgunImpl
}

func NewEmailSvc(cfg *config.EmailConfig) IEmailService {
	mailgunClient := mailgun.NewMailgun(cfg.Domain, cfg.APIKey)

	return &emailServiceImpl{
		client: mailgunClient,
	}
}

func (e *emailServiceImpl) SendEmail(email, subject, body, to string) error {
	logger.Info("Sending email", "to", to, "subject", subject, "body", body)

	message := e.client.NewMessage(email, subject, body, email)

	_, _, err := e.client.Send(message)
	return err
}
