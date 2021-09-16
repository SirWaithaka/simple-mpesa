package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"simple-mpesa/src"
	"simple-mpesa/src/storage"
)

// NewDatabase creates a new Database object
func NewDatabase(config src.Config) (*storage.Database, error) {
	var err error

	// var db *storage.Database
	db := new(storage.Database)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		// DriverName: config.Driver,
		DSN: config.DSN,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open error: %v", err)
	}

	// get native database/sql connection
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("native sql fetch error: %v", err)
	}

	// test connection to db
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping error: %v", err)
	}

	db.DB = gormDB

	return db, nil
}
