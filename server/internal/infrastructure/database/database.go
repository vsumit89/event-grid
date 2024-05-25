package db

import "server/internal/config"

type IDatabase interface {
	Connect(cfg *config.DBConfig) error
	GetClient() interface{}
	Migrate() error
	Close() error
}

func NewDBService(cfg *config.DBConfig) (IDatabase, error) {
	return newPgDB(cfg)
}
