package daos

import (
	"database/sql"
	"time"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/mitchellh/mapstructure"
)

type PostsDAO struct { }

type IPostsDAO interface {
    CreatePost(executor db.QueryExecutor, spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time, username string) (*responses.PostPreview, error) 
    GetPostProperties(executor db.QueryExecutor, postID string, spotifyID string) (*responses.PostPreview, error)
    GetUserPostsProperties(executor db.QueryExecutor, spotifyID string, createdAt time.Time) ([]responses.PostPreview, error)
    GetPostVotes(executor db.QueryExecutor, postID string, spotifyID string) ([]responses.UserIdentifer, []responses.UserIdentifer, error)
    RemovePostVote(executor db.QueryExecutor, voterSpotifyID string, posterSpotifyID string, songID string) error 
    UpdatePost(executor db.QueryExecutor, spotifyID string, songID string, updatePostRequest *requests.UpdatePostRequestDTO, username string) (*responses.PostPreview, error) 
    LikePost(executor db.QueryExecutor, spotifyID string, posterSpotifyID string, songID string) error
    DislikePost(executor db.QueryExecutor, spotifyID string, posterSpotifyID string, songID string) error
    DeletePost(executor db.QueryExecutor, songID string, spotifyID string) error
    GetPostComments(executor db.QueryExecutor, spotifyID string, songID string, paginationKey time.Time) ([]responses.Comment, error)
}

func(p *PostsDAO) CreatePost(executor db.QueryExecutor, spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time, username string) (*responses.PostPreview, error) {

	query := `INSERT INTO posts (albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := executor.Exec(query, albumImage, albumID, albumName, createdAt, rating, songID, songName, text, createdAt, spotifyID)

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

func(p *PostsDAO) GetPostVotes(executor db.QueryExecutor, postID string, spotifyID string)  ([]responses.UserIdentifer, []responses.UserIdentifer, error) {
        query2 := `SELECT post_votes.voterspotifyid, users.username, post_votes.liked 
                   FROM post_votes INNER JOIN users ON post_votes.voterspotifyid = users.spotifyid
                   WHERE post_votes.posterspotifyid = $1 AND post_votes.postsongid = $2`

       likes := []responses.UserIdentifer{}
       dislikes := []responses.UserIdentifer{}


       rows, err := executor.Query(query2, spotifyID, postID)

        if err != nil {
            return nil, nil, customerrors.WrapBasicError(err)
        }

        for rows.Next() {
            userID := responses.UserIdentifer{}
            liked := true
            rows.Scan(&userID.SpotifyID, &userID.Username, &liked)
            if liked {
                likes = append(likes, userID)
            } else {
                dislikes = append(dislikes, userID)
            }
        }

        return likes, dislikes, nil

}
func(p *PostsDAO) GetPostProperties(executor db.QueryExecutor, postID string, spotifyID string) (*responses.PostPreview, error) {
    post := &responses.PostPreview{}

    query := `SELECT albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid, username 
              FROM posts 
              INNER JOIN users ON users.spotifyid = posts.posterspotifyid 
              WHERE posts.posterspotifyid = $1 AND posts.songid = $2`

    row := executor.QueryRow(query, spotifyID, postID)

    albumArtUri := sql.NullString{}

    err := row.Scan(&albumArtUri,
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
        return nil, customerrors.WrapBasicError(err)
    }

    post.AlbumArtURI = albumArtUri.String

    return post, nil

}

func(p *PostsDAO) RemovePostVote(executor db.QueryExecutor, voterSpotifyID string, posterSpotifyID string, songID string) error {
	query := `DELETE FROM post_votes WHERE voterspotifyid = $1 AND posterspotifyid = $2 AND postsongid = $3`

	res, err := executor.Exec(query, voterSpotifyID, posterSpotifyID, songID)

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


func(p *PostsDAO) DeletePost(executor db.QueryExecutor, songID string, spotifyID string) error {
	query := `DELETE FROM posts WHERE posterspotifyid = $1 AND songid = $2`

	res, err := executor.Exec(query, spotifyID, songID)

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

func(p *PostsDAO) UpdatePost(executor db.QueryExecutor, spotifyID string, songID string, updatePostRequest *requests.UpdatePostRequestDTO, username string) (*responses.PostPreview, error) {

    postPreview := &responses.PostPreview{}

    updatedPostRequestMap := make(map[string]any)
    mapstructure.Decode(updatePostRequest, &updatedPostRequestMap)

    t := time.Now().UTC()
    updatedPostRequestMap["updatedat"] = &t

    conditionals := make(map[string]any)
    conditionals["posterspotifyid"] = spotifyID
    conditionals["songid"] = songID

    returning := []string{"albumarturi", "albumid", "albumname", "createdat", "rating", "songid", "songname", "review", "updatedat", "posterspotifyid"}

    query, vals := db.PatchQueryBuilder("posts", updatedPostRequestMap, conditionals, returning)

    res := executor.QueryRow(query, vals...)

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
        return nil, err
    }

	return postPreview, nil
}

func(p *PostsDAO) LikePost(executor db.QueryExecutor, spotifyID string, posterSpotifyID string, songID string) error {

    query := `INSERT INTO post_votes (voterspotifyid, posterspotifyid, postsongid, createdat, updatedat, liked) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              ON CONFLICT (voterspotifyid, posterspotifyid, postsongid) DO UPDATE SET updatedat=$5, liked=$6`

    res, err := executor.Exec(query,
        spotifyID,
        posterSpotifyID,
        songID,
        time.Now().UTC(),
        time.Now().UTC(),
        true)

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

func(p *PostsDAO) DislikePost(executor db.QueryExecutor, spotifyID string, posterSpotifyID string, songID string) error {

    query := `INSERT INTO post_votes (voterspotifyid, posterspotifyid, postsongid, createdat, updatedat, liked) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              ON CONFLICT (voterspotifyid, posterspotifyid, postsongid) DO UPDATE SET updatedat=$5, liked=$6`

    res, err := executor.Exec(query,
        spotifyID,
        posterSpotifyID,
        songID,
        time.Now().UTC(),
        time.Now().UTC(),
        false)

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

func(p *PostsDAO) GetPostComments(executor db.QueryExecutor, spotifyID string, songID string, paginationKey time.Time) ([]responses.Comment, error) {

    query := `SELECT commentid, commentorspotifyid, posterspotifyid, songid, commenttext, createdat, updatedat 
              FROM comments
              WHERE posterspotifyid = $1 AND songid = $2 AND createdAt < $3 
              ORDER BY createdat DESC 
              LIMIT 25 `


    rows, err := executor.Query(query, spotifyID, songID, paginationKey)


    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    comments := []responses.Comment{}

    for rows.Next() {

        comment := &responses.Comment{}
        err := rows.Scan(&comment.CommentID, &comment.CommentorID, &comment.PostSpotifyID, &comment.SongID, &comment.CommentText, &comment.CreatedAt, &comment.UpdatedAt)

        if err != nil {
            return nil, customerrors.WrapBasicError(err)
        }

        comments = append(comments, *comment)

    }


    return comments, nil
}


func(p *PostsDAO) GetUserPostsProperties(executor db.QueryExecutor, spotifyID string, createdAt time.Time) ([]responses.PostPreview, error) {
    query := `SELECT posts.albumarturi, posts.albumid, posts.albumname, posts.createdat, posts.rating, posts.songid, posts.songname, posts.review, posts.updatedat, posts.posterspotifyid, users.username
                FROM posts 
                INNER JOIN users 
                ON users.spotifyid = posts.posterspotifyid
                WHERE posts.posterspotifyid = $1 AND posts.createdat < $2 ORDER BY posts.createdat LIMIT 25 `

    postPreviews := []responses.PostPreview{}

        rows, err := executor.Query(query, spotifyID, createdAt)

        if err != nil {
            return nil, customerrors.WrapBasicError(err)
        }


        for rows.Next() {
            post := responses.PostPreview{}
            albumArtUri := sql.NullString{}
            err := rows.Scan(&albumArtUri, &post.AlbumID, &post.AlbumName, &post.CreatedAt, &post.Rating, &post.SongID, &post.SongName, &post.Text, &post.UpdatedAt, &post.SpotifyID, &post.Username)
            if err != nil {
                return nil, customerrors.WrapBasicError(err)
            }
            post.AlbumArtURI = albumArtUri.String
            post.Likes = []responses.UserIdentifer{}
            post.Dislikes = []responses.UserIdentifer{}
            postPreviews = append(postPreviews, post)
        }

        return postPreviews, nil
}
