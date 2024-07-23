goose_dbstring ?= postgresql://postgres:04122001@localhost:5432/tunes
goose_dir ?= ./db/migrations

goose-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(goose_dbstring) GOOSE_MIGRATION_DIR=$(goose_dir) goose up
goose-reset:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(goose_dbstring) GOOSE_MIGRATION_DIR=$(goose_dir) goose reset


