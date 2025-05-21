package mockshared

import (
	"food-story/shared/config"
	"time"
)

func MockupConfig() config.Config {
	cfg := config.Config{
		AppPort:       "8080",
		AppEnv:        "local",
		AppHost:       "localhost",
		FrontendURL:   "http://localhost:3000",
		SecretKey:     "KDMe9hXvOTas9UzJIE0LOYHjXauakXmj",
		JwtSecret:     "",
		JwtExpireMs:   "",
		RedisAddress:  "localhost:6379",
		RedisPassword: "",
		KafkaBrokers:  "localhost:9092",
		TimeZone:      "Asia/Bangkok",
	}
	cfg.TableSessionDuration = 1 * time.Hour
	return cfg
}
