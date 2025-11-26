package config

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/utils"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type DBConfig struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBDatabase string `mapstructure:"DB_DATABASE"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBSchema   string `mapstructure:"DB_SCHEMA"`
	TimeZone   string `mapstructure:"TZ"`
}

func InitDBConfig(envFile string) (DBConfig, error) {

	var cfg DBConfig

	if _, err := os.Stat(envFile); err == nil {
		viper.SetConfigFile(envFile)
		readCfgErr := viper.ReadInConfig()
		if readCfgErr != nil {
			return DBConfig{}, readCfgErr
		}
	} else {
		viper.SetDefault("DB_HOST", os.Getenv("DB_HOST"))
		viper.SetDefault("DB_PORT", os.Getenv("DB_PORT"))
		viper.SetDefault("DB_DATABASE", os.Getenv("DB_DATABASE"))
		viper.SetDefault("DB_USERNAME", os.Getenv("DB_USERNAME"))
		viper.SetDefault("DB_PASSWORD", os.Getenv("DB_PASSWORD"))
		viper.SetDefault("DB_SCHEMA", os.Getenv("DB_SCHEMA"))
		viper.SetDefault("TZ", os.Getenv("TZ"))
	}

	_ = viper.Unmarshal(&cfg)

	if !utils.IsValidTimeZone(cfg.TimeZone) {
		return DBConfig{}, errors.New("invalid time zone")
	}

	return cfg, nil
}

func (d *DBConfig) ConnectToDatabase() (*pgxpool.Pool, error) {
	source := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s search_path=%s sslmode=disable TimeZone=%s", d.DBUsername, d.DBPassword, d.DBHost, d.DBPort, d.DBDatabase, d.DBSchema, d.TimeZone)
	dbConn, err := pgxpool.New(context.Background(), source)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
