package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUrl string
}

func LoadConfigurations() (Config, error) {

	if err := godotenv.Load(".env"); err != nil {
		return Config{}, err
	}

	var conf Config
	conf.MongoUrl = os.Getenv("mongourl")

	return conf, nil
}
