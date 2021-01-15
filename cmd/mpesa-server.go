package main

import (
	"fmt"
	"log"
	"os"

	"simple-mpesa/src"
	"simple-mpesa/src/configs"
	"simple-mpesa/src/domain"
	"simple-mpesa/src/registry"
	"simple-mpesa/src/routing"
	"simple-mpesa/src/storage/postgres"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// read yaml config file. Dont pass path to read
	// from default path
	yamlConfig := configs.ReadYaml("")
	config := src.GetConfig(*yamlConfig)

	database, err := postgres.NewDatabase(config)
	if err != nil {
		log.Printf("database err %s", err)
		os.Exit(1)
	}

	// run migrations; update tables
	postgres.Migrate(database)

	channels := registry.NewChannels()
	domains := domain.NewDomain(config, database, channels)

	// create the fiber server.
	server := routing.Router(domains, config) // add endpoints

	// listen and serve
	port := fmt.Sprintf(":%v", 6700)
	log.Fatal(server.Listen(port))
}
