package db

import (
	"context"
	"errors"
	"os"

	"github.com/Jack-Gitter/tunes/models"
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

func GetUserFromDbBySpotifyID(spotifyID string) (*models.User, error) {
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
        return nil, errors.New("could not find user in database")
    } 
    properties, _ := res.Records[0].Get("properties")

    user := &models.User{}
    mapstructure.Decode(properties, user)

    return user, nil

}

func InsertUserIntoDB(spotifyID string, username string, role string) (*models.User, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MERGE (u:User {spotifyID: $spotifyID, username: $username, bio: $bio, role: $role}) return properties(u) as properties",
        map[string]any{
            "spotifyID": spotifyID,
            "username": username,
            "role": role,
            "bio": "",
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase("neo4j"),
    )

    if err != nil {
        return nil, err
    }

    if len(resp.Records) < 1 {
        return nil, errors.New("could not find user in database")
    }

    properties, _ := resp.Records[0].Get("properties")

    user := &models.User{}
    mapstructure.Decode(properties, user)

    return user, nil
}



