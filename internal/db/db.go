package db

import (
	"context"
	"strings"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/jack5341/otel-map-server/internal/models"
)

func Open(dsn string, logLevel string) (*gorm.DB, error) {
	gormDB, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel(logLevel)),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	// Ensure connection is alive at startup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := gormDB.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}

	// Skip OtelTrace auto-migration as it conflicts with ClickHouse schema
	err = gormDB.AutoMigrate(&models.SessionToken{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func gormLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}
