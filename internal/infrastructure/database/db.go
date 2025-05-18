package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gps-no-server/internal/common/config"
	"gps-no-server/internal/core/models"
	"log"
	"os"
	"time"
)

type GormDB struct {
	DB *gorm.DB
}

func NewGormDB(cfg *config.DatabaseConfig) (*GormDB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
		})

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
		func() string {
			if cfg.SSLMode {
				return "require"
			}
			return "disable"
		}(),
		cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&models.Station{},
		&models.Cluster{},
		&models.Ranging{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &GormDB{DB: db}, nil
}

func (g *GormDB) Close() error {
	db, err := g.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
