package db

import (
	"errors"
	"os"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/* ===================== CREATE =====================  */

func CreatePost(post *models.Post, spotifyID string) error {
    _, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) 
    MERGE (p:Post {songID: $songID, songName: $songName, albumName: $albumName, albumArtURI: $albumArtURI, albumID: $albumID, rating: $rating, text: $text, timestamp: $timestamp})
     CREATE (u)-[:Posted]->(p)
     RETURN properties(p) `,
        map[string]any{ 
            "songID": post.SongID,
            "songName": post.SongName,
            "albumName": post.AlbumName,
            "albumArtURI": post.AlbumArtURI,
            "albumID": post.AlbumID,
            "rating": post.Rating,
            "text": post.Text,
            "spotifyID": spotifyID,
            "timestamp": post.Timestamp,
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return err
    }

    return nil
}

/* ===================== READ =====================  */

func GetUserPostByID(postID string, spotifyID string) (*models.Post, bool, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) where p.songID = $postID return properties(p) as Post`,
        map[string]any{ 
            "spotifyID": spotifyID,
            "postID": postID,
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, false, err
    }

    if len(resp.Records) < 1 {
        return nil, false, nil 
    }

    postResponse, found := resp.Records[0].Get("Post")

    if !found {
        return nil, false, errors.New("post has no properites in DB, something went wrong")
    }

    user, found, err := GetUserFromDbBySpotifyID(spotifyID)

    if err != nil {
        return nil, false, err
    }

    if !found {
        return nil, false, nil
    }

    post := &models.Post{}
    post.SpotifyID = spotifyID
    post.Username = user.Username
    mapstructure.Decode(postResponse, post)

    return post, true, nil
}

func GetUserPostsPreviewsByUserID(spotifyID string, username string) ([]models.PostPreview, error) {
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

func GetUserPostPreviewByID(songID string, spotifyID string) (*models.PostPreview, bool, error){
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) where p.songID = $postID return properties(p) as Post`,
        map[string]any{ 
            "spotifyID": spotifyID,
            "postID": songID,
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, false, err
    }

    if len(resp.Records) < 1 {
        return nil, false, nil 
    }

    postResponse, found := resp.Records[0].Get("Post")

    if !found {
        return nil, false, errors.New("post has no properites in DB, something went wrong")
    }

    user, found, err := GetUserFromDbBySpotifyID(spotifyID)

    if err != nil {
        return nil, false, err
    }

    if !found {
        return nil, false, nil
    }

    post := &models.PostPreview{}
    post.SpotifyID = spotifyID
    post.Username = user.Username
    mapstructure.Decode(postResponse, post)

    return post, true, nil
}


/* ===================== DELETE =====================  */

func DeletePost(songID string, spotifyID string) (bool, bool, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) where p.songID = $postID DETACH DELETE p return properties(p) as Post`,
        map[string]any{ 
            "spotifyID": spotifyID,
            "postID": songID,
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return false, false, err
    }

    if len(resp.Records) < 1 {
        return false, false, nil
    }

    return true, true, nil
}


/* PROPERTY UPDATES */
func UpdatePostPropertiesByID() {}
