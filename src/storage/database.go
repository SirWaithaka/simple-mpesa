package storage

import (
	"fmt"

	"gorm.io/gorm"
)

// Database is a wrapper type for the gorm DB pointer
type Database struct {
	*gorm.DB
}

func (db *Database) Close() {
	// w
	d, err := db.DB.DB()
	err = d.Close()
	if err != nil {
		fmt.Printf("Error closing db: %v\n", err)
	}
}
