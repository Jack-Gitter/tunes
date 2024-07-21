package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Jack-Gitter/tunes/models/responses"
	_ "github.com/lib/pq"
)

/* ===================== CREATE =====================  */

func CreatePost(spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time) error {
    query := `INSERT INTO posts 
                    (albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid) 
                    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

    _, err := DB.Driver.Exec(query, albumImage, albumID, albumName, createdAt, rating, songID, songName, text, createdAt, spotifyID)

    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    return nil
}

/* ===================== READ =====================  */

func GetUserPostByID(postID string, spotifyID string) (*responses.Post, bool, error) {
    query := `SELECT albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid, username 
                FROM posts INNER JOIN users ON users.spotifyid = posts.posterspotifyid WHERE posts.posterspotifyid = $1 AND posts.songid = $2`

    row := DB.Driver.QueryRow(query, spotifyID, postID)

    post := &responses.Post{}
    err := row.Scan(&post.AlbumArtURI, &post.AlbumID, &post.AlbumName, &post.CreatedAt, &post.Rating, &post.SongID, &post.SongName, &post.Text, &post.UpdatedAt, &post.SpotifyID, &post.Username)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, false, nil 
        } 
        return nil, false, err
    }

    return post, true, nil
}

// make this method get the posts with id offset -> offset+limit-1
func GetUserPostsPreviewsByUserID(spotifyID string, createdAt time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {
    return nil, nil
}

func GetUserPostPreviewByID(songID string, spotifyID string) (*responses.PostPreview, bool, error){
    return nil, false, nil
}


/* ===================== DELETE =====================  */

func DeletePost(songID string, spotifyID string) (bool, bool, error) {
    return false, false, nil
}


/* PROPERTY UPDATES */
func UpdatePost(spotifyID string, songID string, text *string, rating *int) (*responses.PostPreview, bool, error) {
    return nil, false, nil
}
