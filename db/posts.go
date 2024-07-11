package db

import (
	"errors"
	"os"

	"github.com/Jack-Gitter/tunes/models"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetUserPostById(postID string, spotifyID string) (*models.Post, error) {

    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) MATCH (u)-[:Posted]->(p) where p.songID = $postID return properties(p)`,
        map[string]any{ 
            "spotifyID": spotifyID,
            "postID": postID,
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return nil, err
    }

    if len(resp.Records) < 1 {
        return nil, errors.New("no song found")
    }

    post := &models.Post{}
    mapstructure.Decode(resp, post)

    return post, nil

}

func CreatePost(post *models.Post, spotifyID string) error {
    // do we need to validate in here whether or not the songname, etc are real?
    resp, err := neo4j.ExecuteQuery(DB.Ctx, DB.Driver, 
    `MATCH (u:User {spotifyID: $spotifyID}) 
     MERGE (p:Post {songID: $songID, songName: $songName, albumName: $albumName, albumArtURI: $albumArtURI, albumID: $albumID, rating: $rating, text: $text})
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
        }, 
        neo4j.EagerResultTransformer,
        neo4j.ExecuteQueryWithDatabase(os.Getenv("DB_NAME")),
    )

    if err != nil {
        return err
    }

    if len(resp.Records) < 1 {
        return errors.New("could not find user in database")
    }

    return nil
}


