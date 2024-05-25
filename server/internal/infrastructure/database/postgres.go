package db

import (
	"fmt"
	"server/internal/config"
	"server/internal/models"
	"server/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// pgDBImpl implements IDatabase
type pgDBImpl struct {
	db *gorm.DB
}

func newPgDB(cfg *config.DBConfig) (*pgDBImpl, error) {
	pgClient := &pgDBImpl{}

	err := pgClient.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return pgClient, nil
}

// Connect method is used to connect to the postgres database
func (p *pgDBImpl) Connect(cfg *config.DBConfig) error {
	// building the db url
	dbUrl := p.getDbUrl(cfg)

	// connecting to the database
	client, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		return err
	}

	// setting the client
	p.db = client

	logger.Info("successfully connected to the database")

	return nil
}

// getDbUrl method is used to build the db url
func (p *pgDBImpl) getDbUrl(cfg *config.DBConfig) string {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode)
	return dbURL
}

// Migrate method is used to migrate the database
// it creates the table if it doesn't exist and it adds / removes / updates columns if needed
func (p *pgDBImpl) Migrate() error {
	return p.db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.EventAttendees{},
	)
}

// GetClient method is used to get the database client
func (p *pgDBImpl) GetClient() interface{} {
	return p.db
}

// Close method is used to close the database connection
func (p *pgDBImpl) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
