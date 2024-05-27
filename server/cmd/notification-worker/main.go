package main

import (
	"flag"
	"server/internal/config"
	db "server/internal/infrastructure/database"
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

	logger.Info("successfully connected to the database", "database", dbClient.GetClient())
}
