package db

import (
	"errors"
	"os"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

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
        return nil, true, errors.New("user within the database has no properties")
    }

    user := &models.User{}
    mapstructure.Decode(userResponse, user)

    return user, true, nil
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
        if !exists { return nil, errors.New("post has no properties in database") }
        post := &models.PostPreview{}
        mapstructure.Decode(postResponse, post)
        post.UserIdentifer.SpotifyID = spotifyID
        post.UserIdentifer.Username = username
        posts = append(posts, (*post))
    }

    return posts, nil
}

func GetUserFromDbBySpotifyID(spotifyID string) (*models.User, bool, error) {

    user, foundUser, err := getUserProperties(spotifyID)

    if err != nil {
        return nil, false, err
    }

    if !foundUser {
        return nil, false, nil
    }

    posts, err := getUserPostsPreviews(spotifyID, user.Username)

    if err != nil {
        return nil, false, err
    }

    user.Posts = posts

    return user, true, nil

}

func InsertUserIntoDB(user *models.User) error {
    _, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MERGE (u:User {spotifyID: $spotifyID, username: $username, bio: $bio, role: $role})",
        map[string]any{
            "spotifyID": user.SpotifyID,
            "username": user.Username,
            "role": user.Role,
            "bio": "",
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    return err
}
