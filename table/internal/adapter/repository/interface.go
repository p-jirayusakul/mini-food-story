package repository

import (
	database "food-story/shared/database/sqlc"
	"food-story/shared/snowflakeid"
	"food-story/table/config"
)

type TableRepoImplement struct {
	config      config.Config
	repository  database.Store
	snowflakeID snowflakeid.SnowflakeInterface
}

func NewRepo(config config.Config, repository database.Store, snowflakeID snowflakeid.SnowflakeInterface) *TableRepoImplement {
	return &TableRepoImplement{
		config,
		repository,
		snowflakeID,
	}
}
