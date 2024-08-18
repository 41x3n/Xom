package bootstrap

import (
	interfaces "github.com/41x3n/Xom/interface"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func NewEnv() *interfaces.Env {
	env := &interfaces.Env{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found, falling back to OS environment variables.")
	} else {
		// Unmarshal the config file into the env struct
		err = viper.Unmarshal(env)
		if err != nil {
			log.Fatal("Environment can't be loaded from .env file: ", err)
		}
	}

	// Fallback to OS environment variables for any missing fields
	env.AppEnv = getEnv("APP_ENV", env.AppEnv)
	env.DSN = getEnv("DSN", env.DSN)
	env.TelegramBotToken = getEnv("TELEGRAM_BOT_TOKEN", env.TelegramBotToken)
	env.RabbitMQURL = getEnv("RABBITMQ_URL", env.RabbitMQURL)
	env.ContextTimeout = getEnvAsInt("CONTEXT_TIMEOUT", env.ContextTimeout)

	// Example condition based on environment
	if env.AppEnv == "development" {
		log.Println("The App is running in development environment")
	}

	return env
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
