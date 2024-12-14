package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var Config *Env

type Env struct {
	PORT           string `mapstructure:"PORT"`
	DATABASE_URL   string `mapstructure:"DATABASE_URL"`
	DATABASE_NAME  string `mapstructure:"DATABASE_NAME"`
	JWT_SECRET_KEY string `mapstructure:"JWT_SECRET_KEY"`
	MODE           string `mapstructure:"MODE"`
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
		"PORT":           Config.PORT,
		"DATABASE_URL":   Config.DATABASE_URL,
		"DATABASE_NAME":  Config.DATABASE_NAME,
		"JWT_SECRET_KEY": Config.JWT_SECRET_KEY,
		"MODE":           Config.MODE,
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
