include .env
export

sqlc:
	@sqlc generate

.PHONY: sqlc
