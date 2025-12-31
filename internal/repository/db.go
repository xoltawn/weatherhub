package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/xoltawn/weatherhub/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the Postgres connection and runs migrations
func InitDB(dsn string) (*gorm.DB, error) {
	// 1. Open Connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Logs SQL queries for debugging
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 2. Configure Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 3. Run Auto-Migrations
	log.Println("Running database migrations...")
	err = db.AutoMigrate(
		&domain.Weather{},
		// Add other domain entities here (e.g., &domain.User{})
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return db, nil
}
