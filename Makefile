migrations-run:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgresql://postgres:04122001@localhost:5432/tunes goose -dir=./db/migrations up

