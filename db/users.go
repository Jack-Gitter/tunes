package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetUserFromDbBySpotifyID(spotifyID string) (*models.User, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
   //"MATCH (u:User {spotifyID: $spotifyID}) return properties(u) as properties",
   "MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) return properties(p) as posts, properties(u) as user",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, err
    }

    if len(res.Records) < 1 {
        return nil, errors.New("could not find user in database")
    } 


    userResponse, found := res.Records[0].Get("user")

    if !found {
        return nil, errors.New("no properties for inserted user in the database")
    }

    posts := []models.PostInformationForUser{}

    for _, record := range res.Records {
        postResponse, _ := record.Get("posts")
        post := &models.PostInformationForUser{}
        fmt.Println(postResponse)
        mapstructure.Decode(postResponse, post)
        fmt.Println(post)
        posts = append(posts, (*post))
    }

    user := &models.User{}
    mapstructure.Decode(userResponse, user)
    user.Posts = posts

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
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    return err
}
