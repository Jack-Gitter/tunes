include .env

goose-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_DIR) goose up
goose-reset:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_DIR) goose reset
docker-start: 
	docker compose --profile backend up --build
docker-stop:
	docker compose --profile backend down -v
api-start: 
	go run main.go

