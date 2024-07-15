package db

import (
	"errors"
	"os"
	"time"

	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/* ===================== CREATE =====================  */

func CreatePost(spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, timestamp time.Time) (*responses.Post, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) 
    MERGE (p:Post {songID: $songID, songName: $songName, albumName: $albumName, albumArtURI: $albumArtURI, albumID: $albumID, rating: $rating, text: $text, timestamp: $timestamp})
     CREATE (u)-[:Posted]->(p)
     RETURN properties(p) as Post, u.username as Username`,
        map[string]any{ 
            "songID": songID,
            "songName": songName,
            "albumName": albumName,
            "albumArtURI": albumImage,
            "albumID": albumID,
            "rating": rating,
            "text": text,
            "spotifyID": spotifyID,
            "timestamp": timestamp,
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    postResponse := &responses.Post{}

    if err != nil {
        return nil, err
    }

    post, _ := resp.Records[0].Get("Post")
    username, _ := resp.Records[0].Get("Username")
    mapstructure.Decode(post, postResponse)
    postResponse.Username = username.(string)
    postResponse.SpotifyID = spotifyID

    return postResponse, nil
}

/* ===================== READ =====================  */

func GetUserPostByID(postID string, spotifyID string) (*responses.Post, bool, error) {
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) where p.songID = $postID return properties(p) as Post, u.username as Username`,
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

    postResponse, foundPost := resp.Records[0].Get("Post")
    usernameResponse, foundUsername := resp.Records[0].Get("Username")

    if !foundPost || !foundUsername {
        return nil, false, errors.New("post or username has no properites in DB, something went wrong")
    }

    post := &responses.Post{}
    post.SpotifyID = spotifyID
    post.Username = usernameResponse.(string)
    mapstructure.Decode(postResponse, post)

    return post, true, nil
}

func GetUserPostsPreviewsByUserID(spotifyID string) ([]responses.PostPreview, error) {
    res, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    "MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) return properties(p) as postProperties, u.username as Username",
        map[string]any{
            "spotifyID": spotifyID,
        }, neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, err
    }

    posts := []responses.PostPreview{}

    for _, record := range res.Records {
        postResponse, exists := record.Get("postProperties")
        usernameResponse, uexists := record.Get("Username")
        if !exists || !uexists { return nil, errors.New("post has no properties in database") }
        post := &responses.PostPreview{}
        mapstructure.Decode(postResponse, post)
        post.UserIdentifer.SpotifyID = spotifyID
        post.UserIdentifer.Username = usernameResponse.(string)
        posts = append(posts, (*post))
    }

    return posts, nil
}

func GetUserPostPreviewByID(songID string, spotifyID string) (*responses.PostPreview, bool, error){
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) where p.songID = $postID return properties(p) as Post, u.username as Username`,
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
    usernameResponse, ufound := resp.Records[0].Get("Username")

    if !found || !ufound {
        return nil, false, errors.New("post or username has no properites in DB, something went wrong")
    }


    post := &responses.PostPreview{}
    post.SpotifyID = spotifyID
    post.Username = usernameResponse.(string)
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
func UpdatePost(text string, rating int) {}
