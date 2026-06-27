package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Dsn            string
	Port           string
	JWTSecret      string
	JWTExpiryHours int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	return &Config{
		Dsn:       os.Getenv("DSN"),
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		// JWTExpiryHours: os.Getenv("JWT_EXPIRES_HOURS"),
	}
}
