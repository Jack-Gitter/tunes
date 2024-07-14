package db

import (
	"errors"
	"os"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// make it so that there are seperate queries for user properties, user posts, user followers, user following

func getUserProperties(spotifyID string) (*models.User, bool, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u:User {spotifyID: $spotifyID}) RETURN properties(u) as userProperties",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, false, err
    }

    if len(res.Records) < 1 {
        return nil, false, nil
    }

    userResponse, found := res.Records[0].Get("userProperties")

    if !found {
        return nil, false, errors.New("user within the database has no properties")
    }

    user := &models.User{}
    mapstructure.Decode(userResponse, user)

    return user, true, nil
}

func getUserPostsPreviews(spotifyID string, username string) ([]models.PostPreview, bool, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) return properties(p) as postProperties",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, false, err
    }

    posts := []models.PostPreview{}

    for _, record := range res.Records {
        postResponse, exists := record.Get("postProperties")
        if postResponse == nil { continue }
        if !exists { return nil, true, errors.New("wtf post has no properties") }
        post := &models.PostPreview{}
        mapstructure.Decode(postResponse, post)
        post.UserIdentifer.SpotifyID = spotifyID
        post.UserIdentifer.Username = username
        posts = append(posts, (*post))
    }

    return posts, true, nil
}

func GetUserFromDbBySpotifyID(spotifyID string) (*models.User, bool, error) {

    user, foundUser, err := getUserProperties(spotifyID)

    if err != nil {
        return nil, false, err
    }

    if !foundUser {
        return nil, false, nil
    }

    posts, foundPosts, err := getUserPostsPreviews(spotifyID, user.Username)

    if err != nil {
        return nil, false, err
    }

    if !foundPosts {
        return nil, false, nil
    }

    user.Posts = posts

    return user, true, nil

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
