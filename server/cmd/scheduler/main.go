package main

import (
	"flag"
	"os"
	"os/signal"
	"server/internal/commons"
	"server/internal/config"
	db "server/internal/infrastructure/database"
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

	// connecting to the database
	dbClient, err := db.NewDBService(cfg.DB)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	logger.Info("connected to database", "database", dbClient)

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

	notificationWorker := workers.NotificationScheduler{
		Dispatcher: workers.NewEventDispatcher(workerCount, func(event *workers.NotificationEvent) {
			logger.Info("dispatching event", "event", event)
		}),
	}

	notificationWorker.Dispatcher.Start()

	go func() {
		mqClient.Consume(ch, &notificationWorker)
	}()

	forever := make(chan os.Signal, 1)

	signal.Notify(forever, os.Interrupt, syscall.SIGTERM)

	<-forever

	logger.Info("stopping application")
}
