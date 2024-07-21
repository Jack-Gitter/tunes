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
    query := `SELECT albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid, username 
                FROM posts INNER JOIN users ON users.spotifyid = posts.posterspotifyid WHERE posts.posterspotifyid = $1 AND posts.songid = $2`

    row := DB.Driver.QueryRow(query, spotifyID, songID)

    postPreview := &responses.PostPreview{}

    err := row.Scan(&postPreview.AlbumArtURI, 
                    &postPreview.AlbumID, 
                    &postPreview.AlbumName, 
                    &postPreview.CreatedAt, 
                    &postPreview.Rating, 
                    &postPreview.SongID, 
                    &postPreview.SongName, 
                    &postPreview.Text, 
                    &postPreview.UpdatedAt, 
                    &postPreview.SpotifyID, 
                    &postPreview.Username)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, false, nil 
        } 
        return nil, false, err
    }

    return postPreview, true, nil

}


/* ===================== DELETE =====================  */

func DeletePost(songID string, spotifyID string) (bool, error) {
    query := `DELETE FROM posts WHERE posterspotifyid = $1 AND songid = $2`

    res, err := DB.Driver.Exec(query, spotifyID, songID)

    if err != nil {
        return false, err
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return false, err
    }

    if rows < 1 {
        return false, nil
    }

    return true, nil
}


/* PROPERTY UPDATES */
func UpdatePost(spotifyID string, songID string, text *string, rating *int) (*responses.PostPreview, bool, error) {
    query := "UPDATE posts SET "

    val := 1
    vals := []any{}
    if text != nil {
        query += fmt.Sprintf("review = $%d", val)
        val+=1
        vals = append(vals, text)
    }

    if rating != nil {
        if val > 1 {
            query += fmt.Sprintf(", rating = $%d", val)
        } else {
            query += fmt.Sprintf("rating = $%d", val)
        }
        vals = append(vals, rating)
        val+=1
    }

    query += fmt.Sprintf("WHERE posterspotifyid = $%d AND songid = $%d", val, val+1)
    vals = append(vals, spotifyID, songID)

    res, err := DB.Driver.Exec(query, vals)

    if err != nil {
       return nil, false, err 
    }

    num, err := res.RowsAffected()

    if err != nil {
        return nil, false, err
    }

    if num < 1 {
        return nil, false, nil
    }



    return nil, false, nil
}
