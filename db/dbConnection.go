package db

import (
	"database/sql"
    _ "github.com/lib/pq"
	"fmt"
)

var DBConnection *sql.DB = nil

func ConnectToDB() *sql.DB {

    connectionString := "host=localhost port=5432 user=postgres password=04122001 dbname=tunes sslmode=disable"
    db, err := sql.Open("postgres", connectionString)

    DBConnection = db

    if err != nil {
        fmt.Println(err)
    }

    r := DBConnection.Ping()


    if r != nil {
        fmt.Println(r)
        db.Close()
    }

    return db

}

