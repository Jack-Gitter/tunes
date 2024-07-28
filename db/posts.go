package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jack-Gitter/tunes/db/helpers"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	_ "github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
)

/* ===================== CREATE =====================  */

func CreatePost(spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time, username string) (*responses.PostPreview, error) {

	query := `INSERT INTO posts (albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := DB.Driver.Exec(query, albumImage, albumID, albumName, createdAt, rating, songID, songName, text, createdAt, spotifyID)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
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

	return postPreview, nil
}

/* ===================== READ =====================  */

func GetUserPostByID(postID string, spotifyID string) (*responses.Post, error) {


    post := &responses.Post{}

    transaction := func() error {
        tx, err := DB.Driver.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        _, err = tx.Exec(`SET TRANSACTION ISOLATION LEVEL REPEATABLE READ`)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        query := `SELECT albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid, username 
                  FROM posts 
                  INNER JOIN users ON users.spotifyid = posts.posterspotifyid 
                  WHERE posts.posterspotifyid = $1 AND posts.songid = $2`

        row := tx.QueryRow(query, spotifyID, postID)

        post := &responses.Post{}
        albumArtUri := sql.NullString{}

        err = row.Scan(&albumArtUri,
            &post.AlbumID,
            &post.AlbumName,
            &post.CreatedAt,
            &post.Rating,
            &post.SongID,
            &post.SongName,
            &post.Text,
            &post.UpdatedAt,
            &post.SpotifyID,
            &post.Username)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        post.AlbumArtURI = albumArtUri.String
        post.Likes = []responses.UserIdentifer{}
        post.Dislikes = []responses.UserIdentifer{}


        query2 := `SELECT post_votes.voterspotifyid, users.username, post_votes.liked 
                   FROM post_votes INNER JOIN users ON post_votes.voterspotifyid = users.spotifyid
                   WHERE post_votes.posterspotifyid = $1 AND post_votes.postsongid = $2 `

        rows, err := tx.Query(query2, spotifyID, postID)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        for rows.Next() {
            userID := &responses.UserIdentifer{}
            liked := true
            err := rows.Scan(&userID.SpotifyID, &userID.Username, &liked)
            if err != nil {
                return customerrors.WrapBasicError(err)
            }
            if liked {
                post.Likes = append(post.Likes, *userID)
            } else {
                post.Dislikes = append(post.Dislikes, *userID)
            }
        }

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil

    } 

    err := helpers.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        return nil, err
    }

	return post, nil
}

func GetUserPostsPreviewsByUserID(spotifyID string, createdAt time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {

	paginationResponse := &responses.PaginationResponse[[]responses.PostPreview, time.Time]{}
    transaction := func() error {

        tx, err := DB.Driver.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        query := `SELECT spotifyid from USERS WHERE spotifyid = $1`

        res, err := tx.Exec(query, spotifyID)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        count, err := res.RowsAffected()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        if count < 1 {
            return customerrors.WrapBasicError(sql.ErrNoRows)
        }


        query = `SELECT posts.albumarturi, posts.albumid, posts.albumname, posts.createdat, posts.rating, posts.songid, posts.songname, posts.review, posts.updatedat, posts.posterspotifyid, users.username
                FROM posts 
                INNER JOIN users 
                ON users.spotifyid = posts.posterspotifyid
                WHERE posts.posterspotifyid = $1 AND posts.createdat < $2 ORDER BY posts.createdat DESC LIMIT 25 `

        rows, err := tx.Query(query, spotifyID, createdAt)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        postPreviewsResponse := []responses.PostPreview{}

        for rows.Next() {
            post := responses.PostPreview{}
            albumArtUri := sql.NullString{}
            err := rows.Scan(&albumArtUri, &post.AlbumID, &post.AlbumName, &post.CreatedAt, &post.Rating, &post.SongID, &post.SongName, &post.Text, &post.UpdatedAt, &post.SpotifyID, &post.Username)
            if err != nil {
                return customerrors.WrapBasicError(err)
            }
            post.AlbumArtURI = albumArtUri.String
            post.Likes = []responses.UserIdentifer{}
            post.Dislikes = []responses.UserIdentifer{}
            postPreviewsResponse = append(postPreviewsResponse, post)
        }

        query = `SELECT post_votes.voterspotifyid, users.username, post_votes.liked FROM post_votes 
                INNER JOIN users ON users.spotifyid = post_votes.voterspotifyid 
                WHERE post_votes.posterspotifyid = $1 AND post_votes.postsongid = $2`

        for i := 0; i < len(postPreviewsResponse); i++ {

            votes, err := tx.Query(query, postPreviewsResponse[i].SpotifyID, postPreviewsResponse[i].SongID)

            if err != nil {
                return customerrors.WrapBasicError(err)
            }

            for votes.Next() {
                vote := responses.UserIdentifer{}
                liked := true
                err := votes.Scan(&vote.SpotifyID, &vote.Username, &liked)
                if err != nil {
                    return customerrors.WrapBasicError(err)
                }
                if liked {
                    fmt.Println(vote)
                    postPreviewsResponse[i].Likes = append(postPreviewsResponse[i].Likes, vote)
                } else {
                    postPreviewsResponse[i].Dislikes = append(postPreviewsResponse[i].Dislikes, vote)
                }
            }
        }

        paginationResponse.PaginationKey = time.Now().UTC()
        paginationResponse.DataResponse = postPreviewsResponse

        if len(postPreviewsResponse) > 0 {
            lastPost := postPreviewsResponse[len(postPreviewsResponse)-1]
            paginationResponse.PaginationKey = lastPost.CreatedAt
        } 

        if err = tx.Commit(); err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil
    }

    err := helpers.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        return nil, err
    }

	return paginationResponse, nil
}

func RemoveVote(voterSpotifyID string, posterSpotifyID string, songID string) error {
	query := `DELETE FROM post_votes WHERE voterspotifyid = $1 AND posterspotifyid = $2 AND postsongid = $3`

	res, err := DB.Driver.Exec(query, voterSpotifyID, posterSpotifyID, songID)

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	if rows < 1 {
		return customerrors.WrapBasicError(sql.ErrNoRows)
	}

	return nil
}

func GetUserPostPreviewByID(songID string, spotifyID string) (*responses.PostPreview, error) {
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
		return nil, customerrors.WrapBasicError(err)
	}

	return postPreview, nil

}

/* ===================== DELETE =====================  */

func DeletePost(songID string, spotifyID string) error {
	query := `DELETE FROM posts WHERE posterspotifyid = $1 AND songid = $2`

	res, err := DB.Driver.Exec(query, spotifyID, songID)

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	if rows < 1 {
		return customerrors.WrapBasicError(sql.ErrNoRows)
	}

	return nil
}

/* PROPERTY UPDATES */
func UpdatePost(spotifyID string, songID string, updatePostRequest *requests.UpdatePostRequestDTO, username string) (*responses.PostPreview, error) {

    updatedPostRequestMap := make(map[string]any)
    mapstructure.Decode(updatePostRequest, &updatedPostRequestMap)

    conditionals := make(map[string]any)
    conditionals["posterspotifyid"] = spotifyID
    conditionals["songid"] = songID

    returning := []string{"albumarturi", "albumid", "albumname", "createdat", "rating", "songid", "songname", "review", "updatedat", "posterspotifyid"}

    query, vals := helpers.PatchQueryBuilder("posts", updatedPostRequestMap, conditionals, returning)
    fmt.Println(query)

	/*query := "UPDATE posts SET "

	val := 1
	vals := []any{}
	if text != nil {
		query += fmt.Sprintf("review = $%d", val)
		val += 1
		vals = append(vals, text)
	}

	if rating != nil {
		if val > 1 {
			query += fmt.Sprintf(", rating = $%d", val)
		} else {
			query += fmt.Sprintf("rating = $%d", val)
		}
		vals = append(vals, rating)
		val += 1
	}

	query += fmt.Sprintf(`
    WHERE posterspotifyid = $%d AND songid = $%d RETURNING 
    albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid
    `, val, val+1)
	vals = append(vals, spotifyID, songID)*/

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
		return nil, customerrors.WrapBasicError(err)
	}

	return postPreview, nil
}

func LikeOrDislikePost(spotifyID string, posterSpotifyID string, songID string, liked bool) error {

	query := `INSERT INTO post_votes (voterspotifyid, posterspotifyid, postsongid, createdat, updatedat, liked) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              ON CONFLICT (voterspotifyid, posterspotifyid, postsongid) 
              DO UPDATE SET updatedat=$5, liked=$6`

	res, err := DB.Driver.Exec(query,
		spotifyID,
		posterSpotifyID,
		songID,
		time.Now().UTC(),
		time.Now().UTC(),
		liked)

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return customerrors.WrapBasicError(err)
	}

	if rows < 1 {
		return customerrors.WrapBasicError(sql.ErrNoRows)
	}

	return nil
}
