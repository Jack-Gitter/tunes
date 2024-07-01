package db

import (
    "fmt"
    "context"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type dbConnection struct {
    Ctx context.Context
    Driver neo4j.DriverWithContext
}

var DB = &dbConnection{}

func ConnectToDB() {
    DB.Ctx = context.Background()
    dbUri := "neo4j://localhost"
    dbUser := "test"
    dbPassword := "testtest"

    var err error = nil
    DB.Driver, err = neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))

    err = DB.Driver.VerifyConnectivity(DB.Ctx)

    if err != nil {
        panic(err)
    }

    fmt.Println("connection established!")
    
}

func GetUserFromDbBySpotifyID(spotifyID string) *neo4j.EagerResult {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
   "MATCH (u:User {spotifyID: $spotifyID}) return u",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase("neo4j"),
    )

    if err != nil {
        return nil
    }

    return res

}

func InsertUserIntoDB(spotifyID string, username string) {
    neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
   "MERGE (p:Person {spotifyID: $spotifyID, username: $username})",
        map[string]any{
            "spotifyID": spotifyID,
            "username": username,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase("neo4j"),
    )
}



