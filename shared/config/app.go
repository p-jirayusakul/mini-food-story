package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	AppPort       string `mapstructure:"APP_PORT"`
	AppEnv        string `mapstructure:"APP_ENV"`
	AppHost       string `mapstructure:"APP_HOST"`
	FrontendURL   string `mapstructure:"FRONTEND_URL"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	JwtSecret     string `mapstructure:"JWT_SECRET"`
	JwtExpireMs   string `mapstructure:"JWT_EXPIRE_MILLISECOND"`
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	KafkaBrokers  string `mapstructure:"KAFKA_BROKERS"`
}

func InitConfig(envFile string) Config {

	currentDir, _ := os.Getwd()
	envFile = currentDir + "/" + envFile

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
		viper.SetDefault("FRONTEND_URL", os.Getenv("FRONTEND_URL"))
		viper.SetDefault("SECRET_KEY", os.Getenv("SECRET_KEY"))
		viper.SetDefault("JWT_SECRET", os.Getenv("JWT_SECRET"))
		viper.SetDefault("JWT_EXPIRE_MILLISECOND", os.Getenv("JWT_EXPIRE_MILLISECOND"))
		viper.SetDefault("REDIS_ADDRESS", os.Getenv("REDIS_ADDRESS"))
		viper.SetDefault("REDIS_PASSWORD", os.Getenv("REDIS_PASSWORD"))
		viper.SetDefault("KAFKA_BROKERS", os.Getenv("KAFKA_BROKERS"))
	}

	_ = viper.Unmarshal(&cfg)

	return cfg
}
