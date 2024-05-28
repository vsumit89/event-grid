package mocks

import (
	"server/internal/infrastructure/mq"

	"github.com/rabbitmq/amqp091-go"
)

type MockMQ struct {
	MockConnect                  func() error
	MockDeclareQueue             func(queueName string) (*amqp091.Channel, error)
	MockDeclareQueueWithExchange func(exchangeName, queueName string) (*amqp091.Channel, error)
	MockPublish                  func(ch *amqp091.Channel, body []byte) error
	MockPublishWithExchange      func(ch *amqp091.Channel, body []byte, exchangeName string) error
	MockConsume                  func(ch *amqp091.Channel, queueName string, handler mq.IMessageHandler)
	MockClose                    func() error
}

func (m *MockMQ) Connect() error {
	if m.MockConnect != nil {
		return m.MockConnect()
	}

	return nil

}

func (m *MockMQ) DeclareQueue(queueName string) (*amqp091.Channel, error) {
	if m.MockDeclareQueue != nil {
		return m.MockDeclareQueue(queueName)
	}

	return nil, nil
}

func (m *MockMQ) DeclareQueueWithExchange(exchangeName, queueName string) (*amqp091.Channel, error) {
	if m.MockDeclareQueueWithExchange != nil {
		return m.MockDeclareQueueWithExchange(exchangeName, queueName)
	}

	return nil, nil
}

func (m *MockMQ) Publish(ch *amqp091.Channel, body []byte) error {
	if m.MockPublish != nil {
		return m.MockPublish(ch, body)
	}

	return nil
}

func (m *MockMQ) PublishWithExchange(ch *amqp091.Channel, body []byte, exchangeName string) error {
	if m.MockPublishWithExchange != nil {
		return m.MockPublishWithExchange(ch, body, exchangeName)
	}

	return nil
}

func (m *MockMQ) Consume(ch *amqp091.Channel, queueName string, handler mq.IMessageHandler) {
	if m.MockConsume != nil {
		m.MockConsume(ch, queueName, handler)
	}
}

func (m *MockMQ) Close() error {
	if m.MockClose != nil {
		return m.MockClose()
	}

	return nil
}
