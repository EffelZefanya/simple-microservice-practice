package platform

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	RedisAddr      string
	RabbitURL      string
	InventoryAddr  string
	AuthToken      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RabbitURL:     getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		InventoryAddr: getEnv("INVENTORY_GRPC_ADDR", "localhost:50051"),
		AuthToken:     getEnv("GRPC_AUTH_TOKEN", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}