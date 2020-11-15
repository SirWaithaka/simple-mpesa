package postgres

import (
	"log"

	"simple-mpesa/app/models"
	"simple-mpesa/app/statement"
	"simple-mpesa/app/storage"
	"simple-mpesa/app/tariff"
)

// Migrate updates the db with new columns, and tables
func Migrate(database *storage.Database) {
	err := database.DB.AutoMigrate(
		models.Admin{},
		models.Agent{},
		models.Merchant{},
		models.Subscriber{},
		models.Account{},
		models.Transaction{},
		statement.Statement{},
		tariff.Charge{},
	)

	if err != nil {
		log.Println(err)
	}
}
