package db

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	//"github.com/mi"
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

func GetUserFromDbBySpotifyID(spotifyID string) (*User, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
   "MATCH (u:User {spotifyID: $spotifyID}) return properties(u) as properties",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase("neo4j"),
    )

    if err != nil {
        return nil, nil
    }

    properties, exists := res.Records[0].Get("properties")

    user := &User{}
    mapstructure.Decode(properties, user)
    
    if exists == false {
        return nil, nil
    }

    return user, nil

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



type User struct {
    Username string
    SpotifyID string
}
