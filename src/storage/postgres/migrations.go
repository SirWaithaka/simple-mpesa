package postgres

import (
	"log"

	"simple-mpesa/src/repositories/schema"
	"simple-mpesa/src/storage"
)

// Migrate updates the db with new columns, and tables
func Migrate(database *storage.Database) {
	err := database.DB.AutoMigrate(
		schema.Admin{},
		schema.Agent{},
		schema.Merchant{},
		schema.Subscriber{},
		schema.Account{},
		schema.Statement{},
		schema.Statement{},
		schema.Charge{},
	)

	if err != nil {
		log.Println(err)
	}
}
