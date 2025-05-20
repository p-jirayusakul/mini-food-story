include .env
export

sqlc:
	@sqlc generate

mock:
	mockgen -package mock -destination shared/mock/database/store.go food-story/shared/database/sqlc Store
	mockgen -package mock -destination shared/mock/snowflake/snowflake.go food-story/shared/snowflakeid SnowflakeInterface
	mockgen -package mock -destination shared/mock/menu/cache/redis.go food-story/menu-service/internal/adapter/cache RedisTableCacheInterface

.PHONY: sqlc mock
