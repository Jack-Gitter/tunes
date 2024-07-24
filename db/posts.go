package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"github.com/Jack-Gitter/tunes/models/responses"
	pq "github.com/lib/pq"
)

/* ===================== CREATE =====================  */

func CreatePost(spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time, username string) (*responses.PostPreview, bool, error) {
    query := `INSERT INTO posts 
                    (albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid) 
                    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

    _, err := DB.Driver.Exec(query, albumImage, albumID, albumName, createdAt, rating, songID, songName, text, createdAt, spotifyID)
    
    if err != nil {
        err, ok := err.(*pq.Error)
        if !ok || err.Code != "23505"{
            return nil, false, err
        }
        return nil, true, nil
    }
    postPreview := &responses.PostPreview{}
    postPreview.SpotifyID = spotifyID
    postPreview.SongID = songID
    postPreview.SongName = songName
    postPreview.AlbumID = albumID
    postPreview.AlbumName = albumName
    postPreview.AlbumArtURI = albumImage
    postPreview.Rating = rating
    postPreview.Text = text
    postPreview.CreatedAt = createdAt
    postPreview.UpdatedAt = createdAt
    postPreview.Username = username



    return postPreview, false, nil
}

/* ===================== READ =====================  */

func GetUserPostByID(postID string, spotifyID string) (*responses.Post, bool, error) {
    query := `SELECT albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid, username 
                FROM posts INNER JOIN users ON users.spotifyid = posts.posterspotifyid WHERE posts.posterspotifyid = $1 AND posts.songid = $2`

    row := DB.Driver.QueryRow(query, spotifyID, postID)

    post := &responses.Post{}
    albumArtUri := sql.NullString{}
    err := row.Scan(&albumArtUri, &post.AlbumID, &post.AlbumName, &post.CreatedAt, &post.Rating, &post.SongID, &post.SongName, &post.Text, &post.UpdatedAt, &post.SpotifyID, &post.Username)
    post.AlbumArtURI = albumArtUri.String

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, false, nil 
        } 
        return nil, false, err
    }

    return post, true, nil
}

func GetUserPostsPreviewsByUserID(spotifyID string, createdAt time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], bool, error) {
    tx, err := DB.Driver.BeginTx(context.Background(), nil)

    if err != nil {
        return nil, false, err
    }

    defer tx.Rollback()

    query :=  `SELECT spotifyid from USERS WHERE spotifyid = $1`

    row := tx.QueryRow(query, spotifyID)

    rep := "" 
    err = row.Scan(&rep)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, false, nil 
        } 
        return nil, false, err
    }

    
    createdAt = createdAt.Round(time.Microsecond)
    query = ` SELECT posts.albumarturi, posts.albumid, posts.albumname, posts.createdat, posts.rating, posts.songid, posts.songname, posts.review, posts.updatedat, posts.posterspotifyid, users.username
            FROM posts 
            INNER JOIN users 
            ON users.spotifyid = posts.posterspotifyid
            WHERE posts.posterspotifyid = $1 AND posts.createdat < $2 ORDER BY posts.createdat DESC LIMIT 25 `


    rows, err := tx.Query(query, spotifyID, createdAt)

    if err != nil {
        return nil, false, err
    }

    postPreviewsResponse := []responses.PostPreview{}

    for rows.Next() {
        post := responses.PostPreview{}
        albumArtUri := sql.NullString{}
        err := rows.Scan(&albumArtUri, &post.AlbumID, &post.AlbumName, &post.CreatedAt, &post.Rating, &post.SongID, &post.SongName, &post.Text, &post.UpdatedAt, &post.SpotifyID, &post.Username)
        post.AlbumArtURI = albumArtUri.String
        if err != nil {
            return nil, false, err
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

     if err = tx.Commit(); err != nil {
         return nil, false, err
    }

    return paginationResponse, true, nil
}

func GetUserPostPreviewByID(songID string, spotifyID string) (*responses.PostPreview, bool, error){
    query := `SELECT albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid, username 
                FROM posts INNER JOIN users ON users.spotifyid = posts.posterspotifyid WHERE posts.posterspotifyid = $1 AND posts.songid = $2`

    row := DB.Driver.QueryRow(query, spotifyID, songID)

    postPreview := &responses.PostPreview{}
    albumArtUri := sql.NullString{}

    err := row.Scan(&albumArtUri,
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

    postPreview.AlbumArtURI = albumArtUri.String

    if err != nil {
        if err == sql.ErrNoRows {
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
func UpdatePost(spotifyID string, songID string, text *string, rating *int, username string) (*responses.PostPreview, bool, error) {

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

    res := DB.Driver.QueryRow(query, vals...)

    postPreview := &responses.PostPreview{}
    albumArtUri := sql.NullString{}
    err := res.Scan(&albumArtUri,
                    &postPreview.AlbumID, 
                    &postPreview.AlbumName, 
                    &postPreview.CreatedAt, 
                    &postPreview.Rating, 
                    &postPreview.SongID, 
                    &postPreview.SongName, 
                    &postPreview.Text, 
                    &postPreview.UpdatedAt, 
                    &postPreview.SpotifyID)

    postPreview.Username = username
    postPreview.AlbumArtURI = albumArtUri.String

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, false, nil 
        } 
        return nil, false, err
    }

    return postPreview, true, nil
}


func LikeOrDislikePost(spotifyID string, posterSpotifyID string, songID string, liked bool) (bool, error) {
    query := "INSERT INTO post_votes (voterspotifyid, posterspotifyid, postsongid, createdat, updatedat, liked) values ($1, $2, $3, $4, $5, $6) ON CONFLICT (voterspotifyid, posterspotifyid, postsongid) DO UPDATE SET updatedat=$5, liked=$6"
    res, err := DB.Driver.Exec(query, spotifyID, posterSpotifyID, songID, time.Now().UTC(), time.Now().UTC(), liked)

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
