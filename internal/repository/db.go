package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/xoltawn/weatherhub/internal/domain"
	"github.com/xoltawn/weatherhub/pkg/errutil"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the Postgres connection and runs migrations
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return db, nil
}

func MapGormError(err error, context string) error {
	if err == nil {
		return nil
	}

	var finalErr error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		finalErr = domain.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		finalErr = domain.ErrAlreadyExists
	default:
		finalErr = domain.ErrInternal
	}

	return errutil.Wrap(finalErr, context)
}
