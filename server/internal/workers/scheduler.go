package workers

import (
	"encoding/json"
	"fmt"
	"server/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
)

type NotificationEvent struct {
	UnixTimestamp int64 `json:"timestamp"`
	EventID       uint  `json:"event_id"`
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

	n.Dispatcher.AddEvent(&event)

}
