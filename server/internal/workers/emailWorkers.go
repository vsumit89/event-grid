package workers

import (
	"encoding/json"
	"server/internal/repository"
	"server/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
)

type EmailWorker struct {
	repo repository.IEventsRepository
}

func NewEmailWorker(repo repository.IEventsRepository) *EmailWorker {
	return &EmailWorker{
		repo: repo,
	}
}

func (e *EmailWorker) HandleMessage(message interface{}) {

	rabbitMqMsg := message.(amqp091.Delivery)

	var event NotificationEvent

	err := json.Unmarshal(rabbitMqMsg.Body, &event)
	if err != nil {
		logger.Error("event read error", "error", err.Error())
		return
	}

	logger.Info("event read successfully", "event", event)

	if event.Kind != "email" {
		rabbitMqMsg.Ack(false)
		return
	}

	logger.Info("event read successfully", "event", event)

	eventDetails, err := e.repo.GetEventByIDOnly(event.EventID)

	if err != nil {
		logger.Error("event read error", "error", err.Error())
		return
	}

	rabbitMqMsg.Ack(false)

	logger.Info("event read successfully", "event", eventDetails)
}
