package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name string
	Mode string
	URL  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	DSN      string
}

type JWTConfig struct {
	Secret    []byte
	TTL       time.Duration
	TTLInHour int
}

var defaultValues = map[string]string{
	"APP_NAME": "golang-todo",
	"GIN_MODE": "release",
	"APP_URL":  "localhost:8080",
	"PORT":     "8080",

	"DB_HOST":     "127.0.0.1",
	"DB_PORT":     "5432",
	"DB_USERNAME": "postgres",
	"DB_PASSWORD": "postgres",
	"DB_NAME":     "todo_list",

	"JWT_SECRET":      "my-secret-key",
	"JWT_EXP_IN_HOUR": "1",
}

// LoadConfig loads configuration from environment variables
// This should be called explicitly once in main.go
func LoadConfig() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file, using environment variables and defaults")
	}

	cfg := &Config{
		App: AppConfig{
			Name: getEnvAsString("APP_NAME"),
			Mode: getEnvAsString("GIN_MODE"),
			URL:  getEnvAsString("APP_URL"),
			Port: getEnvAsString("PORT"),
		},
		Database: DatabaseConfig{
			Host:     getEnvAsString("DB_HOST"),
			Port:     getEnvAsInt("DB_PORT"),
			Username: getEnvAsString("DB_USERNAME"),
			Password: getEnvAsString("DB_PASSWORD"),
			Name:     getEnvAsString("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:    []byte(getEnvAsString("JWT_SECRET")),
			TTLInHour: getEnvAsInt("JWT_EXP_IN_HOUR"),
			TTL:       time.Duration(getEnvAsInt("JWT_EXP_IN_HOUR")) * time.Hour,
		},
	}

	// Build database DSN
	cfg.Database.DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	return cfg, nil
}

func getEnvAsString(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if def, ok := defaultValues[key]; ok {
		return def
	}

	return ""
}

func getEnvAsInt(key string) int {
	valueStr := getEnvAsString(key)
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("%s not an integer", key)

		if def, ok := defaultValues[key]; ok {
			if defInt, err := strconv.Atoi(def); err == nil {
				return defInt
			}
		}
		return 0
	}

	return valueInt
}
