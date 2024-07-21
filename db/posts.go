package db

import (
	"context"
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
    query := `
            SELECT posts.albumarturi, posts.albumid, posts.albumname, posts.createdat, posts.rating, posts.songid, posts.songname, posts.review, posts.updatedat, posts.posterspotifyid, users.username
            FROM posts 
            INNER JOIN users 
            ON users.spotifyid = posts.posterspotifyid
            WHERE posts.posterspotifyid = $1 AND posts.createdat < $2 ORDER BY posts.createdat DESC LIMIT 25 `

    rows, err := DB.Driver.Query(query, spotifyID, createdAt)

    postPreviewsResponse := []responses.PostPreview{}

    for rows.Next() {
        post := responses.PostPreview{}
        err := rows.Scan(&post.AlbumArtURI, &post.AlbumID, &post.AlbumName, &post.CreatedAt, &post.Rating, &post.SongID, &post.SongName, &post.Text, &post.UpdatedAt, &post.SpotifyID, &post.Username)
        if err != nil {
            return nil, err
        }
        postPreviewsResponse = append(postPreviewsResponse, post)
    }

    paginationResponse := &responses.PaginationResponse[[]responses.PostPreview, time.Time]{}
    paginationResponse.DataResponse = postPreviewsResponse

    if len(postPreviewsResponse) > 0 {
        lastPost := postPreviewsResponse[len(postPreviewsResponse)-1]
        paginationResponse.PaginationKey = lastPost.CreatedAt
    } else {
        paginationResponse.PaginationKey = time.Now().UTC()
    }

    if err != nil {
        return nil, err
    }

    return paginationResponse, nil
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
    tx, err := DB.Driver.BeginTx(context.Background(), nil)

    if err != nil {
        return nil, false, err
    }

    defer tx.Rollback()

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

    query += fmt.Sprintf(`
    WHERE posterspotifyid = $%d AND songid = $%d RETURNING 
    albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid
    `, val, val+1)
    vals = append(vals, spotifyID, songID)

    res := tx.QueryRow(query, vals...)


    postPreview := &responses.PostPreview{}
    err = res.Scan(&postPreview.AlbumArtURI, &postPreview.AlbumID, &postPreview.AlbumName, &postPreview.CreatedAt, &postPreview.Rating, &postPreview.SongID, &postPreview.SongName, 
&postPreview.Text, &postPreview.UpdatedAt, &postPreview.SpotifyID)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, false, nil 
        } 
        return nil, false, err
    }


    query = "SELECT username FROM users WHERE spotifyid = $1"
    res = tx.QueryRow(query, spotifyID)

    err = res.Scan(&postPreview.Username)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, false, nil 
        } 
        return nil, false, err
    }

     if err = tx.Commit(); err != nil {
         return nil, false, err
    }

    return postPreview, true, nil
}
