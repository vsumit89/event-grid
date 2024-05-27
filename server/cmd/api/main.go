package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"server/internal/commons"
	"server/internal/config"
	"server/internal/handlers"
	db "server/internal/infrastructure/database"
	"server/internal/infrastructure/mq"
	"server/internal/infrastructure/transport"
	"server/internal/repository"
	"server/internal/services"
	"server/pkg/logger"
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

	// connecting to the database
	dbClient, err := db.NewDBService(cfg.DB)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	logger.Info("successfully connected to the database")

	logger.Info("running migrations")
	// migrating the database
	err = dbClient.Migrate()
	if err != nil {
		logger.Error("error while migrating the data", "error", err.Error())
		return
	}

	mqClient := mq.NewMessageQueue(cfg.Queue)
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	err = mqClient.Connect()
	if err != nil {
		logger.Error("error while starting application", "error", err.Error())
		return
	}

	logger.Info("successfully connected to the message queue")

	logger.Info("successfully migrated the database")

	jwtSvc := commons.NewJWTService(cfg.JWT)

	// instantiating the repositories, services and handlers
	userRepoSvc := repository.NewUserRepository(dbClient)

	eventRepoSvc := repository.NewEventsRepository(dbClient)

	userSvc := services.NewUserSvc(&services.UserSvcOptions{
		Repository: userRepoSvc,
	})

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: eventRepoSvc,
		UserRepository:  userRepoSvc,
		MQ:              mqClient,
	})

	// deps contains the dependencies for the handlers
	deps := handlers.Container{
		UserSvc:  userSvc,
		JWTSvc:   jwtSvc,
		EventSvc: eventSvc,
	}

	// getting all the handlers
	router := handlers.GetRoutes(&deps)

	httpServer := transport.NewHTTPServer(&cfg.Server, router)
	if err != nil {
		logger.Error("error while http server", "error", err.Error())
		return
	}

	terminationSignal := make(chan os.Signal, 1)

	signal.Notify(terminationSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		// starting the http server
		logger.Info("starting http server")
		err = httpServer.Run()
		if err != nil {
			logger.Error("error while starting http server", "error", err.Error())
			return
		}

		logger.Info("successfully started http server", "port", cfg.Server.Port)
	}()

	<-terminationSignal

	logger.Info("starting graceful shutdown")

	err = dbClient.Close()
	if err != nil {
		logger.Error("error while closing the database connection", "error", err.Error())
		return
	}

	logger.Info("successfully closed the database connection")

	// stopping the http server
	err = httpServer.Stop()
	if err != nil {
		logger.Error("error while stopping the http server", "error", err.Error())
		return
	}

	logger.Info("graceful shutdown completed")
}
