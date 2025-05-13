package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // For loading .env file.
)

// Holds the application configuration.
type Config struct {
	APIGatewayPort           string
	//UserProductServiceGRPCHost string
}

// Gets the .env values or returns a default one.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Loads configuration from environment variables or .env file.
func LoadConfig() *Config {
	// Load .env file if it exists.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading from environment variables.")
	}

	return &Config{
		APIGatewayPort:           getEnv("API_GATEWAY_PORT", "8080"),
		//UserProductServiceGRPCHost: getEnv("USER_PRODUCT_SERVICE_GRPC_HOST", "localhost:50051"),
	}
}

