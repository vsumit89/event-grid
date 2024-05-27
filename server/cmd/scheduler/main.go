package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"server/internal/commons"
	"server/internal/config"
	"server/internal/infrastructure/mq"
	"server/internal/workers"
	"syscall"

	"server/pkg/logger"
	"strings"
)

func main() {
	// reading the l (log level) flag from the command line
	logLevel := flag.String("l", "INFO", "log level")

	// initializing the logger with the provided options
	logger.InitLogger(
		logger.WithLogFmt(logger.TEXT),
		logger.WithLevel(logger.LOG_LEVEL(strings.ToUpper(*logLevel))),
	)

	// reading the c (config path) flag from the command line
	configPath := flag.String("c", "./config.yaml", "path for reading config")

	// parsing the flags
	flag.Parse()

	logger.Info("starting application")

	var err error

	workerCount := 4

	// reading the configuration
	cfg, err := config.InitConfig(*configPath)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	mqClient := mq.NewMessageQueue(cfg.Queue)

	err = mqClient.Connect()
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	ch, err := mqClient.DeclareQueue(commons.QueueName)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	emailCh, err := mqClient.DeclareQueue(commons.EmailQueue)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	notificationWorker := workers.NotificationScheduler{
		Dispatcher: workers.NewEventDispatcher(workerCount, func(event *workers.NotificationEvent) {
			logger.Info("dispatching event", "event", event)

			event.Kind = "email"

			jsonBody, err := json.Marshal(event)
			if err != nil {
				logger.Error("error while dispatching event", "error", err.Error())
				return
			}

			err = mqClient.Publish(emailCh, jsonBody)
			if err != nil {
				logger.Error("error while dispatching event", "error", err.Error())
				return
			}
		}),
	}

	notificationWorker.Dispatcher.Start()

	go func() {
		mqClient.Consume(ch, commons.QueueName, &notificationWorker)
	}()

	forever := make(chan os.Signal, 1)

	signal.Notify(forever, os.Interrupt, syscall.SIGTERM)

	<-forever

	logger.Info("stopping application")

	err = mqClient.Close()
	if err != nil {
		logger.Error("error while stopping application", "error", err.Error())
		return
	}
}
