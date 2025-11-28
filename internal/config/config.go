package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                       int
	Env                        string
	DBDriver                   string
	DBSource                   string
	JWTSecret                  string
	JWTAccessExpirationMinutes time.Duration
	JWTRefreshExpirationDays   time.Duration
}

// AppConfig holds the loaded configuration
var AppConfig *Config

// LoadConfig loads environment variables from .env file
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	port, _ := strconv.Atoi(getEnv("PORT", "3000"))
	accessExp, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRATION_MINUTES", "30"))
	refreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION_DAYS", "30"))

	AppConfig = &Config{
		Port:                       port,
		Env:                        getEnv("NODE_ENV", "development"),
		DBDriver:                   getEnv("DB_DRIVER", "sqlite"),
		DBSource:                   getEnv("DB_SOURCE", "app.db"),
		JWTSecret:                  getEnv("JWT_SECRET", "secret"),
		JWTAccessExpirationMinutes: time.Duration(accessExp) * time.Minute,
		JWTRefreshExpirationDays:   time.Duration(refreshExp) * 24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}