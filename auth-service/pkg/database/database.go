package database

import (
	"auth-service/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	Connect(cfg config.ConfigProvider) error
	AutoMigrate(models ...interface{}) error
	GetDB() (readDB *gorm.DB, writeDB *gorm.DB)
	GetReadDB() *gorm.DB
	GetWriteDB() *gorm.DB
	Close() error
}

type GormDatabase struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func (g *GormDatabase) Connect(cfg config.ConfigProvider) error {
	// replica
	dsnRead := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.GetDBReadHost(), cfg.GetDBReadPort(), cfg.GetDBReadUser(), cfg.GetDBReadPassword(), cfg.GetDBReadName(), cfg.GetDBSSLMode())

	readDB, err := gorm.Open(postgres.Open(dsnRead), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to read database: %w", err)
	}

	sqlReadDB, err := readDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get read database instance: %w", err)
	}
	sqlReadDB.SetMaxOpenConns(25)
	sqlReadDB.SetMaxIdleConns(10)
	sqlReadDB.SetConnMaxLifetime(5 * time.Minute)

	// primary
	dsnWrite := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.GetDBWriteHost(), cfg.GetDBWritePort(), cfg.GetDBWriteUser(), cfg.GetDBWritePassword(), cfg.GetDBWriteName(), cfg.GetDBSSLMode())

	writeDB, err := gorm.Open(postgres.Open(dsnWrite), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to write database: %w", err)
	}

	sqlWriteDB, err := writeDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get write database instance: %w", err)
	}
	sqlWriteDB.SetMaxOpenConns(25)
	sqlWriteDB.SetMaxIdleConns(10)
	sqlWriteDB.SetConnMaxLifetime(5 * time.Minute)

	g.writeDB = writeDB
	g.readDB = readDB

	return nil
}

func (g *GormDatabase) AutoMigrate(models ...interface{}) error {
	if len(models) > 0 {
		if err := g.writeDB.AutoMigrate(models...); err != nil {
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}
	return nil
}

func (g *GormDatabase) GetDB() (readDB *gorm.DB, writeDB *gorm.DB) {
	return g.readDB, g.writeDB
}

func (g *GormDatabase) GetReadDB() *gorm.DB {
	return g.readDB
}

func (g *GormDatabase) GetWriteDB() *gorm.DB {
	return g.writeDB
}

func (g *GormDatabase) Close() error {
	if g.writeDB != nil {
		sqlDB, err := g.writeDB.DB()
		if err != nil {
			return fmt.Errorf("failed to get write database instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	if g.readDB != nil {
		sqlDB, err := g.readDB.DB()
		if err != nil {
			return fmt.Errorf("failed to get read database instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	return nil
}
