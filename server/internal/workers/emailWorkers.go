package workers

import (
	"encoding/json"
	"fmt"
	"server/internal/infrastructure/email"
	"server/internal/models"
	"server/internal/repository"
	"server/pkg/logger"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type EmailWorker struct {
	repo     repository.IEventsRepository
	userRepo repository.IUserRepository
	emailSvc email.IEmailService
}

func NewEmailWorker(repo repository.IEventsRepository, user repository.IUserRepository, emailSvc email.IEmailService) *EmailWorker {
	return &EmailWorker{
		repo:     repo,
		emailSvc: emailSvc,
		userRepo: user,
	}
}

func (e *EmailWorker) HandleMessage(message interface{}) {

	rabbitMqMsg := message.(amqp091.Delivery)

	var event NotificationEvent

	err := json.Unmarshal(rabbitMqMsg.Body, &event)
	if err != nil {
		logger.Error("event read error", "error", err.Error())
		rabbitMqMsg.Nack(false, false)
		return
	}

	logger.Info("event read successfully", "event", event)

	logger.Info("event read successfully", "event", event)

	eventDetails, err := e.repo.GetEventByIDOnly(event.EventID)
	if err != nil {
		logger.Error("event read error", "error", err.Error())
		rabbitMqMsg.Nack(false, false)
		return
	}

	userDetails, err := e.userRepo.GetUserByID(event.CreatedBy)
	if err != nil {
		logger.Error("user read error", "error", err.Error())
		rabbitMqMsg.Nack(false, false)
		return
	}

	for _, attendee := range eventDetails.Attendees {
		text := generatePlainTextEmail(*eventDetails, attendee)

		err = e.emailSvc.SendEmail(userDetails.Email, "Invitation for the Event", text, attendee.Email)
		if err != nil {
			logger.Error("email sending error", "error", err.Error())
			rabbitMqMsg.Nack(false, false)
			return
		}
	}

	rabbitMqMsg.Ack(false)

	logger.Info("event read successfully", "event", eventDetails)
}

func generatePlainTextEmail(event models.Event, attendee models.User) string {
	offset := 5*time.Hour + 30*time.Minute
	return fmt.Sprintf(`Hi %s,

You are invited to the following event:

Event: %s
Description: %s
Date: %s
Start Time: %s
End Time: %s
Location: %s

We hope you can join us!

Best regards,

`,
		attendee.Name,
		event.Title,
		event.Description,
		event.Start.Add(offset).Format("January 2, 2006"),
		event.Start.Add(offset).Format("3:04 PM"),
		event.End.Add(offset).Format("3:04 PM"),
		event.MeetingURL,
	)
}
