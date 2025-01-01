package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var Config *Env

type Env struct {
	PORT          string `mapstructure:"PORT"`
	DATABASE_URL  string `mapstructure:"DATABASE_URL"`
	DATABASE_NAME string `mapstructure:"DATABASE_NAME"`

	JWT_TOKEN_COOKIE    string `mapstructure:"JWT_TOKEN_COOKIE"`
	JWT_DECODED_PAYLOAD string `mapstructure:"JWT_DECODED_PAYLOAD"`
	JWT_SECRET_KEY      string `mapstructure:"JWT_SECRET_KEY"`

	// Email Configuration
	SMTP_HOST     string `mapstructure:"SMTP_HOST"`
	SMTP_PORT     string `mapstructure:"SMTP_PORT"`
	SMTP_USERNAME string `mapstructure:"SMTP_USERNAME"`
	SMTP_PASSWORD string `mapstructure:"SMTP_PASSWORD"`
	FROM_EMAIL    string `mapstructure:"FROM_EMAIL"`

	// RabbitMQ Configuration
	RABBITMQ_URL string `mapstructure:"RABBITMQ_URL"`
	QUEUE_NAME   string `mapstructure:"QUEUE_NAME"`

	// Redis Configuration
	REDIS_ENDPOINT string `mapstructure:"REDIS_ENDPOINT"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	REDIS_PORT     string `mapstructure:"REDIS_PORT"`

	// Mode for golang
	MODE string `mapstructure:"MODE"`
}

func NewEnv() error {
	path, _ := os.Getwd()
	os.Chdir(path + "/../config")
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("can't find the environment file : %v", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("environment can't be loaded: %v", err)
	}

	envMap := map[string]string{
		"PORT":             Config.PORT,
		"DATABASE_URL":     Config.DATABASE_URL,
		"DATABASE_NAME":    Config.DATABASE_NAME,
		"JWT_TOKEN_COOKIE": Config.JWT_TOKEN_COOKIE,
		"JWT_SECRET_KEY":   Config.JWT_SECRET_KEY,
		"SMTP_HOST":        Config.SMTP_HOST,
		"SMTP_PORT":        Config.SMTP_PORT,
		"SMTP_USERNAME":    Config.SMTP_USERNAME,
		"SMTP_PASSWORD":    Config.SMTP_PASSWORD,
		"FROM_EMAIL":       Config.FROM_EMAIL,
		"RABBITMQ_URL":     Config.RABBITMQ_URL,
		"QUEUE_NAME":       Config.QUEUE_NAME,
		"REDIS_ENDPOINT":   Config.REDIS_ENDPOINT,
		"REDIS_PASSWORD":   Config.REDIS_PASSWORD,
		"REDIS_PORT":       Config.REDIS_PORT,
		"MODE":             Config.MODE,
	}

	for key, value := range envMap {
		if value == "" {
			return fmt.Errorf("environment variable %s is not set", key)
		}
	}

	if Config.MODE == "development" {
		log.Println("Server started in development mode")
	}

	return nil
}
