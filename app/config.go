package app

import (
	"fmt"

	"simple-wallet/configs"
)

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func (d Database) String(sslmode string) string {
	return fmt.Sprintf(""+
		"host=%s "+
		"port=%s "+
		"user=%s "+
		"dbname=%s "+
		"password=%s "+
		"sslmode=%v"+
		"", d.Host, d.Port, d.User, d.DBName, d.Password, sslmode)
}

type Config struct {
	DB Database

	Secret string
}

func GetConfig(cfg configs.YamlConfig) Config {
	return Config{
		DB: Database{
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			DBName:   cfg.Database.DBName,
		},

		Secret: cfg.AppSecret,
	}
}
