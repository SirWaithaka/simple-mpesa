package main

import (
	"fmt"
	"log"
	"os"

	"simple-mpesa/app"
	"simple-mpesa/app/registry"
	"simple-mpesa/app/routing"
	"simple-mpesa/app/storage/postgres"
	"simple-mpesa/configs"
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
