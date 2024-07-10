package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// make it so that there are seperate queries for user properties, user posts, user followers, user following

func getUserProperties(spotifyID string) (*models.User, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u:User {spotifyID: $spotifyID}) RETURN properties(u) as userProperties",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, err
    }

    if len(res.Records) < 1 {
        return nil, errors.New("could not find user with that ID in the DB")
    }

    userResponse, found := res.Records[0].Get("userProperties")

    if !found {
        return nil, errors.New("user has no properties in the DB, wtf?")
    }

    user := &models.User{}
    mapstructure.Decode(userResponse, user)

    return user, nil
}

func getUserPostsPreviews(spotifyID string, username string) ([]models.PostPreview, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) return properties(p) as postProperties",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, err
    }

    posts := []models.PostPreview{}

    for _, record := range res.Records {
        postResponse, exists := record.Get("postProperties")
        if postResponse == nil { continue }
        if !exists { return nil, errors.New("wtf post has no properties") }
        post := &models.PostPreview{}
        mapstructure.Decode(postResponse, post)
        post.UserIdentifer.SpotifyID = spotifyID
        post.UserIdentifer.Username = username
        posts = append(posts, (*post))
    }
    return posts, nil
}
func GetUserFromDbBySpotifyID(spotifyID string) (*models.User, error) {

    user, err := getUserProperties(spotifyID)

    if err != nil {
        return nil, err
    }

    posts, err := getUserPostsPreviews(spotifyID, user.Username)

    if err != nil {
        fmt.Println("second")
        return nil, err
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
/*func GetUserFromDbBySpotifyID(spotifyID string) (*models.User, error) {
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

}*/

