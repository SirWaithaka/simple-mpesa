package postgres

import (
	"log"

	"simple-mpesa/app/models"
	"simple-mpesa/app/storage"
)

// Migrate updates the db with new columns, and tables
func Migrate(database *storage.Database) {
	err := database.DB.AutoMigrate(
		models.Admin{},
		models.Agent{},
		models.Merchant{},
		models.Subscriber{},
		models.User{},
		models.Account{},
		models.Transaction{},
	)

	if err != nil {
		log.Println(err)
	}
}
