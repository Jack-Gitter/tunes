package db

import (
	"errors"
	"os"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func DoesUserPostExist(postID string, spotifyID string) (bool, error) {

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
        return false, err
    }

    if len(resp.Records) < 1 {
        return false, err
    }

    _, found := resp.Records[0].Get("Post")

    if !found {
        return true, errors.New("post has no properites in DB, something went wrong")
    }

    return true, nil

}

func GetUserPostByID(postID string, spotifyID string) (*models.Post, bool, error) {

    user, found, err := GetUserFromDbBySpotifyID(spotifyID)

    if err != nil {
        return nil, false, err
    }

    if !found {
        return nil, false, nil
    }

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

    post := &models.Post{}
    post.SpotifyID = spotifyID
    post.Username = user.Username
    mapstructure.Decode(postResponse, post)

    return post, true, nil

}

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

func DeletePost(songID string, spotifyID string) (bool, bool, error) {

    found, err := DoesUserPostExist(songID, spotifyID)
    if err != nil {
        return false, false, err
    }

    if !found {
        return false, false, nil
    }

    _, err = neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) where p.songID = $postID DETACH DELETE p`,
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

    return true, true, nil


}
