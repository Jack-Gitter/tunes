goose_dbstring := postgresql://postgres:04122001@localhost:5432/tunes 
migrations-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgresql://postgres:04122001@localhost:5432/tunes GOOSE_MIGRATION_DIR=./db/migrations goose up
migrations-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgresql://postgres:04122001@localhost:5432/tunes GOOSE_MIGRATION_DIR=./db/migrations goose reset


