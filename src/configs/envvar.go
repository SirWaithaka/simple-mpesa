package configs

import "github.com/kelseyhightower/envconfig"

type EnvVarConfig struct {
	PORT      string `envconfig:"port" default:"6700"`
	DBDriver  string `envconfig:"db_driver"`
	DBDSN     string `envconfig:"db_dsn"`
	SecretKey string `envconfig:"secret_key"`
}

func GetEnvConfig() (*EnvVarConfig, error) {
	var cfg EnvVarConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
