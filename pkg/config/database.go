package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type DBConfig struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBDatabase string `mapstructure:"DB_DATABASE"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBSchema   string `mapstructure:"DB_SCHEMA"`
}

func InitDBConfig(envFile string) DBConfig {

	var cfg DBConfig

	if _, err := os.Stat(envFile); err == nil {
		viper.SetConfigFile(envFile)
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	} else {
		viper.SetDefault("DB_HOST", os.Getenv("DB_HOST"))
		viper.SetDefault("DB_PORT", os.Getenv("DB_PORT"))
		viper.SetDefault("DB_DATABASE", os.Getenv("DB_DATABASE"))
		viper.SetDefault("DB_USERNAME", os.Getenv("DB_USERNAME"))
		viper.SetDefault("DB_PASSWORD", os.Getenv("DB_PASSWORD"))
		viper.SetDefault("DB_SCHEMA", os.Getenv("DB_SCHEMA"))
	}

	_ = viper.Unmarshal(&cfg)

	return cfg
}
