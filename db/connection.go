package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type dbConnection struct {
	Driver *sql.DB
}

type QueryExecutor interface {
    Exec(query string, args ...interface{}) (sql.Result, error) // if we have to swap out databases, sql.Result is a big problem? -- no, its a driver over multiple dbs!
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
}

func ConnectToDB() *sql.DB {
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

    return db

}
