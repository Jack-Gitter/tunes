package db

import (
	"database/sql"
    _ "github.com/lib/pq"
	"fmt"
)

func ConnectToDB() *sql.DB {
    connectionString := "host=localhost port=5432 user=postgres password=04122001 dbname=tunes sslmode=disable"
    db, err := sql.Open("postgres", connectionString)

    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("db success!")
    return db

}

