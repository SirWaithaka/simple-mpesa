package postgres

import (
	"log"

	"simple-wallet/app/models"
	"simple-wallet/app/storage"
)

// Migrate updates the db with new columns, and tables
func Migrate(database *storage.Database) {
	err := database.DB.AutoMigrate(
		models.User{},
		models.Account{},
		models.Transaction{},
	)

	if err != nil {
		log.Println(err)
	}
}
