package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"simple-mpesa/src"
	"simple-mpesa/src/configs"
	"simple-mpesa/src/domain"
	"simple-mpesa/src/registry"
	"simple-mpesa/src/routing"
	"simple-mpesa/src/storage/postgres"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// load .env file from current directory into env variables
	err := godotenv.Load(".env")
	if err != nil {
		// we will just log the error and ignore then continue
		log.Printf("error reading .env file: %v", err)
	}

	// load env vars into struct
	envconfig, err := configs.GetEnvConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	// transform envvar config into app config
	config := src.GetAppConfig(envconfig)

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
