package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type Redis struct {
	Address  string
	Password string
}

type Config struct {
	Port          string
	DatabaseURL   string
	Redis         Redis
	SecretKey     string
	FrontendURL   string
	ResendAPIKey  string
	EmailFrom     string
	EmailFromName string
}

func LoadConfig() (Config, error) {
	var config Config

	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			return config, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	config.Port = mustGetEnv("PORT")
	config.SecretKey = mustGetEnv("AUTH_SECRET_KEY")

	config.DatabaseURL = mustGetEnv("DB_URL")

	config.Redis = Redis{
		Address:  mustGetEnv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	config.FrontendURL = mustGetEnv("FRONTEND_URL")
	config.ResendAPIKey = mustGetEnv("RESEND_API_KEY")
	config.EmailFrom = mustGetEnv("EMAIL_FROM")
	config.EmailFromName = mustGetEnv("EMAIL_FROM_NAME")

	if _, err := strconv.Atoi(config.Port); err != nil {
		return config, fmt.Errorf("invalid port number: %w", err)
	}

	return config, nil
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}
