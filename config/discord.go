package config

import (
	"os"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	TOKEN string
}

func LoadConfig() (BotConfig, error) {
	err := godotenv.Load(".env")
	cfg := BotConfig{
		TOKEN: os.Getenv("TOKEN"),
	}
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
