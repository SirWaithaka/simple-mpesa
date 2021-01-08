package postgres

import (
	"log"

	"simple-mpesa/src/account"
	"simple-mpesa/src/admin"
	"simple-mpesa/src/agent"
	"simple-mpesa/src/merchant"
	"simple-mpesa/src/storage"
	"simple-mpesa/src/subscriber"
	"simple-mpesa/src/tariff"
	"simple-mpesa/src/transaction"
)

// Migrate updates the db with new columns, and tables
func Migrate(database *storage.Database) {
	err := database.DB.AutoMigrate(
		admin.Administrator{},
		agent.Agent{},
		merchant.Merchant{},
		subscriber.Subscriber{},
		account.Account{},
		transaction.Statement{},
		account.Statement{},
		tariff.Charge{},
	)

	if err != nil {
		log.Println(err)
	}
}
