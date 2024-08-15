package bootstrap

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

func NewEnv(depth int) *Env {
	env := Env{}

	var path []string
	for i := 1; i < depth; i++ {
		path = append(path, "..")
	}

	projectRoot, err := filepath.Abs(filepath.Join(path...)) // Move up two directories
	if err != nil {
		log.Fatalf("Error getting project root path: %v", err)
		return nil
	}

	fmt.Printf("ENV path: %v", path)

	viper.SetConfigFile(filepath.Join(projectRoot, ".env"))

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Cant find the .env file: ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment couldn't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The app is running in development evn")
	}

	return &env
}
