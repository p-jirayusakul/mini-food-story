package repository

import (
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/snowflakeid"
)

type Implement struct {
	config      config.Config
	repository  database.Store
	snowflakeID snowflakeid.SnowflakeInterface
}

func NewRepo(config config.Config, repository database.Store, snowflakeID snowflakeid.SnowflakeInterface) *Implement {
	return &Implement{
		config,
		repository,
		snowflakeID,
	}
}
