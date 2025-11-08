package database

import (
	"github.com/Alfian57/golang-todo/pkg/config"
	"github.com/Alfian57/golang-todo/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates a new database connection with the given configuration
// This function should be called once in main.go and the db instance should be passed to repositories
func New(cfg *config.Config, log logger.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database",
			logger.F("host", cfg.Database.Host),
			logger.F("database", cfg.Database.Name),
			logger.F("port", cfg.Database.Port),
			logger.F("error", err),
		)
		return nil, err
	}

	log.Info("Successfully connected to database",
		logger.F("host", cfg.Database.Host),
		logger.F("database", cfg.Database.Name),
		logger.F("port", cfg.Database.Port),
	)

	return db, nil
}
