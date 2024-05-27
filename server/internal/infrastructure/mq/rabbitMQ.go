package mq

import (
	"server/internal/commons"
	"server/internal/config"
	"server/pkg/logger"

	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type rabbitMQClient struct {
	client *amqp091.Connection
	queue  amqp091.Queue
	cfg    *config.QueueConfig
}

func newRabbitMQClient(cfg *config.QueueConfig) *rabbitMQClient {
	return &rabbitMQClient{
		cfg: cfg,
	}
}

func (r *rabbitMQClient) Connect() error {
	logger.Info("connecting to the queue")

	var err error
	url := fmt.Sprintf("%s://%s:%s@%s:%s/", r.cfg.Protocol, r.cfg.Username, r.cfg.Password, r.cfg.Host, r.cfg.Port)

	i := 0
	for {
		r.client, err = amqp091.Dial(url)
		if err != nil {
			if i >= 3 {
				return err
			}
		} else {
			break
		}

		time.Sleep(5 * time.Second)
		i++
	}

	logger.Info("connected to the queue")
	return nil
}

func (r *rabbitMQClient) DeclareQueue(queueName string) (*amqp091.Channel, error) {
	var err error

	ch, err := r.client.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		commons.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	r.queue, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments

	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queueName,
		"notification",
		commons.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (r *rabbitMQClient) DeclareQueueWithExchange(exchangeName, queueName string) (*amqp091.Channel, error) {
	var err error

	ch, err := r.client.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	r.queue, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments

	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queueName,
		"notification",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (r *rabbitMQClient) Publish(ch *amqp091.Channel, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		commons.ExchangeName,
		"notification",
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logger.Info("Error publishing message", "error", err)
		return err
	}

	return nil
}

func (r *rabbitMQClient) Consume(ch *amqp091.Channel, queueName string, handler IMessageHandler) {
	logger.Info("consumer started")

	logger.Info("binding queue", "queue", r.queue.Name)

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		logger.Info("Error consuming message", "error", err)
		return
	}

	// forever channel makes the runner alive indefinitely
	var forever chan struct{}

	go func() {
		for d := range msgs {
			go handler.HandleMessage(d)
		}
	}()

	<-forever
}

func (r *rabbitMQClient) Close() error {
	return r.client.Close()
}
