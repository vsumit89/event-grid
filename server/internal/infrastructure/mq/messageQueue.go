package mq

import (
	"server/internal/config"

	"github.com/rabbitmq/amqp091-go"
)

type IMessageQueue interface {
	Connect() error
	DeclareQueue(queueName string) (*amqp091.Channel, error)
	DeclareQueueWithExchange(exchangeName, queueName string) (*amqp091.Channel, error)
	Publish(ch *amqp091.Channel, body []byte) error
	PublishWithExchange(ch *amqp091.Channel, body []byte, exchangeName string) error
	Consume(ch *amqp091.Channel, queueName string, handler IMessageHandler)
	Close() error
}

type IMessageHandler interface {
	HandleMessage(interface{})
}

func NewMessageQueue(cfg *config.QueueConfig) IMessageQueue {
	return newRabbitMQClient(cfg)
}
