package workers

import (
	"encoding/json"
	"fmt"
	"server/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
)

type NotificationEvent struct {
	UnixTimestamp int64  `json:"timestamp"`
	EventID       uint   `json:"event_id"`
	Kind          string `json:"kind"`
	CreatedBy     uint   `json:"created_by"`
}

type NotificationScheduler struct {
	Dispatcher *EventDispatcher
}

func NewNotificationWorker(dispatcher *EventDispatcher) *NotificationScheduler {
	return &NotificationScheduler{
		Dispatcher: dispatcher,
	}
}

func (n *NotificationScheduler) HandleMessage(message interface{}) {
	rabbitMqMsg := message.(amqp091.Delivery)

	var event NotificationEvent

	err := json.Unmarshal(rabbitMqMsg.Body, &event)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.Info("event read successfully", "event", event)

	if event.Kind == "update" {
		logger.Info("calling remove event", "event", event.EventID)

		n.Dispatcher.RemoveEvent(event.EventID)
	}

	n.Dispatcher.AddEvent(&event)

	rabbitMqMsg.Ack(false)
}
