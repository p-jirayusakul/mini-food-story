package repository

import (
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/snowflakeid"
)

type ProductRepoImplement struct {
	config      config.Config
	repository  database.Store
	snowflakeID snowflakeid.SnowflakeInterface
}

func NewRepo(config config.Config, repository database.Store, snowflakeID snowflakeid.SnowflakeInterface) *ProductRepoImplement {
	return &ProductRepoImplement{
		config,
		repository,
		snowflakeID,
	}
}
