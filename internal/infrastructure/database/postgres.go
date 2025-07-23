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

type PostgresConnection struct {
	db *gorm.DB
}

func NewPostgresConnection(cfg *config.DatabaseConfig) (*PostgresConnection, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
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

	return &PostgresConnection{db: db}, nil
}

func (p *PostgresConnection) GetDB() *gorm.DB {
	return p.db
}

func (p *PostgresConnection) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresConnection) Migrate() error {
	return p.db.AutoMigrate(
		&models.Station{},
		&models.Cluster{},
		//&models.Ranging{},
		//&models.StationConfiguration{},
	)
}
