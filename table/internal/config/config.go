package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	AppPort     string `mapstructure:"APP_PORT"`
	AppEnv      string `mapstructure:"APP_ENV"`
	AppHost     string `mapstructure:"APP_HOST"`
	SecretKey   string `mapstructure:"SECRET_KEY"`
	JwtSecret   string `mapstructure:"JWT_SECRET"`
	JwtExpireMs string `mapstructure:"JWT_EXPIRE_MILLISECOND"`
}

func InitConfig(envFile string) Config {

	currentDir, _ := os.Getwd()
	envFile = currentDir + "/" + envFile
	fmt.Println("Current Directory:", envFile)

	var cfg Config

	if _, err := os.Stat(envFile); err == nil {
		viper.SetConfigFile(envFile)
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	} else {
		fmt.Println("Config file not found")
		viper.SetDefault("APP_PORT", os.Getenv("APP_PORT"))
		viper.SetDefault("APP_ENV", os.Getenv("APP_ENV"))
		viper.SetDefault("APP_HOST", os.Getenv("APP_HOST"))
		viper.SetDefault("SECRET_KEY", os.Getenv("SECRET_KEY"))
		viper.SetDefault("JWT_SECRET", os.Getenv("JWT_SECRET"))
		viper.SetDefault("JWT_EXPIRE_MILLISECOND", os.Getenv("JWT_EXPIRE_MILLISECOND"))
	}

	_ = viper.Unmarshal(&cfg)

	return cfg
}
