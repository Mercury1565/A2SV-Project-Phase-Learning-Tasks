package bootstrap

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                string `mapstructure:"APP_ENV"`
	ServerAddress         string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout        int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                string `mapstructure:"DB_HOST"`
	DBPort                string `mapstructure:"DB_PORT"`
	DBName                string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret     string `mapstructure:"ACCESS_TOKEN_SECRET"`
}

func NewEnv() *Env {
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Println("Failed to load .env file, falling back to system environment variables")
	}

	viper.AutomaticEnv() // read from environment variables

	env := &Env{
		ServerAddress:         viper.GetString("SERVER_ADDRESS"),
		AppEnv:                viper.GetString("APP_ENV"),
		ContextTimeout:        viper.GetInt("CONTEXT_TIMEOUT"),
		DBHost:                viper.GetString("DB_HOST"),
		DBPort:                viper.GetString("DB_PORT"),
		DBName:                viper.GetString("DB_NAME"),
		AccessTokenExpiryHour: viper.GetInt("ACCESS_TOKEN_EXPIRY_HOUR"),
		AccessTokenSecret:     viper.GetString("ACCESS_TOKEN_SECRET"),
	}

	if env.ServerAddress == "" {
		log.Fatal("SERVER_ADDRESS not set")
	}

	if env.AppEnv == "development" {
		log.Println("The app is running in development env")
	}

	return env
}
