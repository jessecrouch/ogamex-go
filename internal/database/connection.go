package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	appLogger "ogamex-go/internal/logger"
	"ogamex-go/internal/schema"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(host string, port int, user, password, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	appLogger.Info().Msg("Running auto-migration...")
	return d.DB.AutoMigrate(
		&schema.User{},
		&schema.Planet{},
		&schema.UserTech{},
		&schema.FleetMission{},
		&schema.BuildingQueue{},
		&schema.ResearchQueue{},
		&schema.UnitQueue{},
	)
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Ping() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *Database) Begin(ctx context.Context) *gorm.DB {
	tx := d.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		appLogger.Warn().Err(tx.Error).Msg("Failed to begin transaction")
	}
	return tx
}

func (d *Database) WithContext(ctx context.Context) *gorm.DB {
	return d.DB.WithContext(ctx)
}
