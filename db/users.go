package db

import (
	"errors"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

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

    properties, found := res.Records[0].Get("properties")

    if !found {
        return nil, errors.New("no properties for inserted user in the database")
    }

    user := &models.User{}
    mapstructure.Decode(properties, user)

    return user, nil

}

func InsertUserIntoDB(spotifyID string, username string, role string) error {
    _, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MERGE (u:User {spotifyID: $spotifyID, username: $username, bio: $bio, role: $role}) return properties(u) as properties",
        map[string]any{
            "spotifyID": spotifyID,
            "username": username,
            "role": role,
            "bio": "",
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase("neo4j"),
    )

    return err
}
