package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/redrru/fantasy-dota/pkg/log"
)

type Config struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type DB struct {
	Gorm *gorm.DB
	sql  *sql.DB
}

func NewDB(config Config) (*DB, error) {
	log.GetLogger().Info(context.Background(), "[DB] New", zap.Any("config", config))

	_, err := pgx.ParseConfig(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse DSN failed: %w", err)
	}

	orm, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{
		Logger: newGormLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("initialize db session failed: %w", err)
	}

	if err := orm.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("use otel plugin failed: %w", err)
	}

	sqlDB, err := orm.DB()
	if err != nil {
		return nil, fmt.Errorf("get generic db object failed: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return &DB{Gorm: orm, sql: sqlDB}, nil
}

func (db *DB) Close() error {
	return db.sql.Close()
}

func (db *DB) Ping(context context.Context) error {
	return db.sql.PingContext(context)
}
