include .env
export

sqlc:
	@sqlc generate

mock:
	mockgen -package mockdb -destination shared/mock/database/store.go food-story/shared/database/sqlc Store
	mockgen -package mockshared -destination shared/mock/shared/snowflake.go food-story/shared/snowflakeid SnowflakeInterface

	mockgen -package mockcache -destination shared/mock/cache/menu/redis.go food-story/menu-service/internal/adapter/cache RedisTableCacheInterface
	mockgen -package mockcache -destination shared/mock/cache/table/redis.go food-story/table-service/internal/adapter/cache RedisTableCacheInterface
	mockgen -package mockcache -destination shared/mock/cache/order/redis.go food-story/order-service/internal/adapter/cache RedisTableCacheInterface

	mockgen -package mockqueue -destination shared/mock/queue/order/queue.go food-story/order-service/internal/adapter/queue/producer QueueProducerInterface


.PHONY: sqlc mock
