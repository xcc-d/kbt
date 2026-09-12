package client

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbConfig struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
}

func NewMysqlClient(cfg DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB: %w", err)
	}
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetConnMaxLifetime(cfg.MaxConnectionLifeTime)
	return gormDB, nil
}
