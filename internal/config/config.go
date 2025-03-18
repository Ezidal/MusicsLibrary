package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	Environment    string
	ExternalApiUrl string
}

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Print("No .env file found")
// 	}
// }

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, set default values")
	}

	config := &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "user"),
		DBPass:         getEnv("DB_PASSWORD", "user"),
		DBName:         getEnv("DB_NAME", "db"),
		Environment:    getEnv("ENVIRONMENT", "local"),
		ExternalApiUrl: getEnv("EXTERNAL_API_URL", "http://localhost:8081") + "/info",
	}
	return config
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
