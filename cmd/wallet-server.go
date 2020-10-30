package main

import (
	"fmt"
	"log"
	"os"

	"simple-wallet/app"
	"simple-wallet/app/registry"
	"simple-wallet/app/routing"
	"simple-wallet/app/storage/postgres"
	"simple-wallet/configs"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// read yaml config file. Dont pass path to read
	// from default path
	cfg := configs.ReadYaml("")
	confg := app.GetConfig(*cfg)

	database, err := postgres.NewDatabase(confg)
	if err != nil {
		log.Printf("database err %s", err)
		os.Exit(1)
	}

	// run migrations; update tables
	postgres.Migrate(database)

	channels := registry.NewChannels()
	domain := registry.NewDomain(confg, database, channels)

	// create the fiber server.
	server := routing.Router(domain, confg) // add endpoints

	// listen and serve
	port := fmt.Sprintf(":%v", 6700)
	log.Fatal(server.Listen(port))
}
