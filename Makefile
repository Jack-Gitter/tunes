include .env

goose-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_DIR) goose up
goose-reset:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_DIR) goose reset
docker-start: 
	docker compose --profile postgres-full up
docker-stop:
	docker compose --profile postgres-full down

