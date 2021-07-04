package src

import (
	"simple-mpesa/src/configs"
)

type Config struct {
	Driver string
	DSN    string

	Secret string

	HTTPPort string
}

func GetAppConfig(cfg *configs.EnvVarConfig) Config {
	return Config{
		Driver:   cfg.DBDriver,
		Secret:   cfg.SecretKey,
		DSN:      cfg.DBDSN,
		HTTPPort: cfg.PORT,
	}
}
