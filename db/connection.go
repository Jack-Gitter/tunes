package db

import (
	"context"
	"os"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type dbConnection struct {
    Ctx context.Context
    Driver neo4j.DriverWithContext
}

var DB = &dbConnection{}

func ConnectToDB() {
    DB.Ctx = context.Background()
    dbUri := os.Getenv("DB_URI")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASS")

    var err error = nil
    DB.Driver, err = neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))

    err = DB.Driver.VerifyConnectivity(DB.Ctx)

    if err != nil {
        panic(err)
    }
}




