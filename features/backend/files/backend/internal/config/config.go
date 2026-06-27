package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	FrontendBaseURL string
	Env             string
}

func Load() (*Config, error) {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		Port:            port,
		FrontendBaseURL: os.Getenv("FRONTEND_BASE_URL"),
		Env:             env,
	}, nil
}
