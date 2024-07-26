package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type dbConnection struct {
	Ctx    context.Context
	Driver *sql.DB
}

var DB = &dbConnection{}

func ConnectToDB() {
	DB.Ctx = context.Background()
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	DB.Driver = db

}
