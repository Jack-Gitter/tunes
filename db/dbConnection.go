package db

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
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
        return nil, err
    }

    if len(res.Records) < 1 {
        return nil, errors.New("err")
    } 
    properties, exists := res.Records[0].Get("properties")

    if exists == false {
        return nil, errors.New("missinng")
    }

    user := &User{}
    mapstructure.Decode(properties, user)
    

    return user, nil

}

func InsertUserIntoDB(spotifyID string, username string, role string) (*User, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MERGE (u:User {spotifyID: $spotifyID, username: $username, role: $role}) return properties(u) as properties",
        map[string]any{
            "spotifyID": spotifyID,
            "username": username,
            "role": role,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase("neo4j"),
    )

    if err != nil {
        return nil, err
    }

    if len(resp.Records) < 1 {
        return nil, errors.New("no properties")
    }
    properties, exists := resp.Records[0].Get("properties")

    if !exists {
        return nil, nil
    }

    user := &User{}
    mapstructure.Decode(properties, user)

    return user, nil
}



