package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token  string `envconfig:"TOKEN" required:"true"`
	Domain string `envconfig:"DOMAIN" required:"true"`
}

func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) && !os.IsNotExist(err) {
		return config, fmt.Errorf("godotenv.load: %w", err)
	}

	err = envconfig.Process("", &config)

	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}
	return config, nil
}
