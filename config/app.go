package config

import (
	"os"
)

type AppConfig struct {
    Port     string
}

func App() *AppConfig {
    return &AppConfig{
		Port:	os.Getenv("APP_PORT"),
	}
}