package db

import (
	"errors"
	"os"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// make it so that there are seperate queries for user properties, user posts, user followers, user following

func GetUserFromDbBySpotifyID(spotifyID string) (*models.User, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "OPTIONAL MATCH (u:User {spotifyID: $spotifyID}) OPTIONAL MATCH (u)-[:Posted]->(p) return properties(p) as post, properties(u) as user",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, err
    }

    userResponse, _ := res.Records[0].Get("user")

    if userResponse == nil {
        return nil, errors.New("could not find user in database")
    } 

    user := &models.User{}
    mapstructure.Decode(userResponse, user)

    posts := []models.PostPreview{}

    for _, record := range res.Records {
        postResponse, _ := record.Get("post")
        if postResponse == nil { continue }
        post := &models.PostPreview{}
        mapstructure.Decode(postResponse, post)
        post.UserIdentifer.SpotifyID = user.SpotifyID
        post.UserIdentifer.Username = user.Username
        posts = append(posts, (*post))
    }

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
