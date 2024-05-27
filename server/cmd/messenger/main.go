package main

import (
	"flag"
	"os"
	"os/signal"
	"server/internal/commons"
	"server/internal/config"
	db "server/internal/infrastructure/database"
	"server/internal/infrastructure/mq"
	"server/internal/repository"
	"server/internal/workers"
	"server/pkg/logger"
	"strings"
	"syscall"
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

	// reading the configuration
	cfg, err := config.InitConfig(*configPath)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	dbClient, err := db.NewDBService(cfg.DB)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	logger.Info("successfully connected to the database")

	mqClient := mq.NewMessageQueue(cfg.Queue)

	err = mqClient.Connect()
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	ch, err := mqClient.DeclareQueueWithExchange(commons.EmailExchange, commons.EmailQueue)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	eventRepo := repository.NewEventsRepository(dbClient)

	emailWorker := workers.NewEmailWorker(eventRepo)

	go func() {
		mqClient.Consume(ch, commons.EmailQueue, emailWorker)
	}()

	forever := make(chan os.Signal, 1)

	signal.Notify(forever, os.Interrupt, syscall.SIGTERM)

	<-forever

}
