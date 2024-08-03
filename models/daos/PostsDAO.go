package daos

import (
	"context"
	"database/sql"
	"net/http"
	"sort"
	"time"

	"github.com/Jack-Gitter/tunes/db"
	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/mitchellh/mapstructure"
)

type PostsDAO struct {
    DB *sql.DB
}

type IPostsDAO interface {
    CreatePost(spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time, username string) (*responses.PostPreview, error) 
    GetUserPostByID(postID string, spotifyID string) (*responses.PostPreview, error) 
    GetUserPostsPreviewsByUserID(spotifyID string, createdAt time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error)
    RemoveVote(voterSpotifyID string, posterSpotifyID string, songID string) error 
    GetUserPostPreviewByID(songID string, spotifyID string) (*responses.PostPreview, error) 
    DeletePost(songID string, spotifyID string) error
    UpdatePost(spotifyID string, songID string, updatePostRequest *requests.UpdatePostRequestDTO, username string) (*responses.PostPreview, error) 
    LikeOrDislikePost(spotifyID string, posterSpotifyID string, songID string, liked bool) error
    GetPostCommentsPaginated(spotifyID string, songID string, paginationKey time.Time) (*responses.PaginationResponse[[]responses.Comment, time.Time], error)
    GetCurrentUserFeed(spotifyID string, t time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error)
}

func(p *PostsDAO) CreatePost(spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time, username string) (*responses.PostPreview, error) {

	query := `INSERT INTO posts (albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := p.DB.Exec(query, albumImage, albumID, albumName, createdAt, rating, songID, songName, text, createdAt, spotifyID)

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

func(p *PostsDAO) GetUserPostByID(postID string, spotifyID string) (*responses.PostPreview, error) {

    post := &responses.PostPreview{}

    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

        if err != nil {
            return err
        }

        err = p.getUserPostProperties(tx, spotifyID, postID, post)

        if err != nil {
            return err
        }

        err = p.getPostLikesAndDislikes(tx, spotifyID, postID, post)

        if err != nil {
            return err
        }

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil

    } 

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        return nil, err
    }

	return post, nil
}

func(p *PostsDAO) GetUserPostsPreviewsByUserID(spotifyID string, createdAt time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {

	paginationResponse := &responses.PaginationResponse[[]responses.PostPreview, time.Time]{PaginationKey: time.Now().UTC()}

    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

        if err != nil {
            return err
        }

        err = p.getUserPostPreviews(tx, spotifyID, createdAt, paginationResponse)

        if err != nil {
            return err
        }

        if err = tx.Commit(); err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil
    }

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        return nil, err
    }

	return paginationResponse, nil
}




func(p *PostsDAO) RemoveVote(voterSpotifyID string, posterSpotifyID string, songID string) error {
	query := `DELETE FROM post_votes WHERE voterspotifyid = $1 AND posterspotifyid = $2 AND postsongid = $3`

	res, err := p.DB.Exec(query, voterSpotifyID, posterSpotifyID, songID)

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

func(p *PostsDAO) GetUserPostPreviewByID(songID string, spotifyID string) (*responses.PostPreview, error) {
	query := `SELECT albumarturi, albumid, albumname, createdat, rating, songid, songname, review, updatedat, posterspotifyid, username 
                FROM posts INNER JOIN users ON users.spotifyid = posts.posterspotifyid WHERE posts.posterspotifyid = $1 AND posts.songid = $2`

	row := p.DB.QueryRow(query, spotifyID, songID)

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

func(p *PostsDAO) DeletePost(songID string, spotifyID string) error {
	query := `DELETE FROM posts WHERE posterspotifyid = $1 AND songid = $2`

	res, err := p.DB.Exec(query, spotifyID, songID)

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

func(p *PostsDAO) UpdatePost(spotifyID string, songID string, updatePostRequest *requests.UpdatePostRequestDTO, username string) (*responses.PostPreview, error) {

    postPreview := &responses.PostPreview{}

    transaction := func() error {
        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

        if err != nil {
            return err
        }

        updatedPostRequestMap := make(map[string]any)
        mapstructure.Decode(updatePostRequest, &updatedPostRequestMap)

        t := time.Now().UTC()
        updatedPostRequestMap["updatedat"] = &t

        conditionals := make(map[string]any)
        conditionals["posterspotifyid"] = spotifyID
        conditionals["songid"] = songID

        returning := []string{"albumarturi", "albumid", "albumname", "createdat", "rating", "songid", "songname", "review", "updatedat", "posterspotifyid"}

        query, vals := db.PatchQueryBuilder("posts", updatedPostRequestMap, conditionals, returning)

        res := tx.QueryRow(query, vals...)

        albumArtUri := sql.NullString{}
        err = res.Scan(&albumArtUri,
            &postPreview.AlbumID,
            &postPreview.AlbumName,
            &postPreview.CreatedAt,
            &postPreview.Rating,
            &postPreview.SongID,
            &postPreview.SongName,
            &postPreview.Text,
            &postPreview.UpdatedAt,
            &postPreview.SpotifyID)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        postPreview.Username = username
        postPreview.AlbumArtURI = albumArtUri.String

        err = p.getPostLikesAndDislikes(tx, spotifyID, songID, postPreview)

        if err != nil {
            return err
        }

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil

    }

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        return nil, err
    }

	return postPreview, nil
}

func(p *PostsDAO) LikeOrDislikePost(spotifyID string, posterSpotifyID string, songID string, liked bool) error {

    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }


        defer tx.Rollback()

        err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

        if err != nil {
            return err
        }

        query := `SELECT COUNT(*)
                  FROM post_votes
                  WHERE voterspotifyid = $1 AND posterspotifyid = $2 AND postsongid = $3 AND liked = $4`

        row := tx.QueryRow(query, spotifyID, posterSpotifyID, songID, liked)

        count := 0
        err = row.Scan(&count)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        if count >= 1 {
            return &customerrors.CustomError{StatusCode: http.StatusConflict, Msg: "cannot like or dislike a post many times"}
        }

        query = `INSERT INTO post_votes (voterspotifyid, posterspotifyid, postsongid, createdat, updatedat, liked) 
                  VALUES ($1, $2, $3, $4, $5, $6) 
                  ON CONFLICT (voterspotifyid, posterspotifyid, postsongid) DO UPDATE SET updatedat=$5, liked=$6`

        res, err := tx.Exec(query,
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

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }
        return nil
    }

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        return err
    }

	return nil
}

func(p *PostsDAO) GetPostCommentsPaginated(spotifyID string, songID string, paginationKey time.Time) (*responses.PaginationResponse[[]responses.Comment, time.Time], error) {

    paginationResponse := &responses.PaginationResponse[[]responses.Comment, time.Time]{PaginationKey: time.Now().UTC()}

    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

        if err != nil {
            return err
        }

        query := `SELECT commentid, commentorspotifyid, posterspotifyid, songid, commenttext, createdat, updatedat 
                  FROM comments
                  WHERE posterspotifyid = $1 AND songid = $2 AND createdAt < $3 
                  ORDER BY createdat DESC 
                  LIMIT 25 `


        rows, err := tx.Query(query, spotifyID, songID, paginationKey)


        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        comments := []responses.Comment{}

        for rows.Next() {

            comment := &responses.Comment{}
            err := rows.Scan(&comment.CommentID, &comment.CommentorID, &comment.PostSpotifyID, &comment.SongID, &comment.CommentText, &comment.CreatedAt, &comment.UpdatedAt)

            if err != nil {
                return customerrors.WrapBasicError(err)
            }

            comments = append(comments, *comment)

        }

        for i := 0; i < len(comments); i++ {

            query = `SELECT liked FROM comment_votes WHERE commentid = $1`
            rows, err := tx.Query(query, comments[i].CommentID)

            if err != nil {
                return customerrors.WrapBasicError(err)
            }

            likes := 0
            dislikes := 0
            for rows.Next() {
               val := true 
               err := rows.Scan(&val)
               if err != nil {
                   return customerrors.WrapBasicError(err)
               }
               if val {
                   likes +=1
               } else {
                   dislikes+=1
               }
            }
               comments[i].Likes = likes
               comments[i].Dislikes = dislikes
        }

        paginationResponse.DataResponse = comments

        if len(comments) > 0 {
            paginationResponse.PaginationKey = comments[len(comments)-1].CreatedAt
        }

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil
    }

    db.RunTransactionWithExponentialBackoff(transaction, 5) 


    return paginationResponse, nil
}

func(p *PostsDAO) GetCurrentUserFeed(spotifyID string, t time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {

    paginationResponse := &responses.PaginationResponse[[]responses.PostPreview, time.Time]{PaginationKey: time.Now().UTC()}

    transaction := func() error {

        tx, err := p.DB.BeginTx(context.Background(), nil)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        defer tx.Rollback()

        query := `SELECT userfollowed FROM followers WHERE follower = $1`

        followedIds := []string{}
        rows, err := p.DB.Query(query, spotifyID)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        for rows.Next() {
            follower := ""
            err := rows.Scan(&follower)
            if err != nil {
                return customerrors.WrapBasicError(err)
            }
            followedIds = append(followedIds, follower)
        }

        posts := []responses.PostPreview{}

        for _, id := range followedIds {
            userPosts := &responses.PaginationResponse[[]responses.PostPreview, time.Time]{}
            err = p.getUserPostPreviews(p.DB, id, t, userPosts)
            if err != nil {
                return err
            }
            posts = append(posts, userPosts.DataResponse...)
        }

        
        sort.Slice(posts, func(i, j int) bool {
            return posts[i].CreatedAt.Before(posts[j].CreatedAt)
        })

        if len(posts) == 0 {
            return nil
        }

        upTo := 0
        if len(posts) > 15 {
            upTo = 15
        }

        paginationResponse.DataResponse = posts[:upTo]
        paginationResponse.PaginationKey = posts[:upTo][upTo-1].CreatedAt

        err = tx.Commit()

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        return nil
    }

    err := db.RunTransactionWithExponentialBackoff(transaction, 5)

    if err != nil {
        return nil, err
    }

    return paginationResponse, nil

}

func(p *PostsDAO) getUserPostProperties(executor db.QueryExecutor, spotifyID string, postID string, post *responses.PostPreview) error {
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
            return customerrors.WrapBasicError(err)
        }

        post.AlbumArtURI = albumArtUri.String

        return nil
}

func(p *PostsDAO) getPostLikesAndDislikes(executor db.QueryExecutor, spotifyID string, postID string, post *responses.PostPreview) error {
        query2 := `SELECT post_votes.voterspotifyid, users.username, post_votes.liked 
                   FROM post_votes INNER JOIN users ON post_votes.voterspotifyid = users.spotifyid
                   WHERE post_votes.posterspotifyid = $1 AND post_votes.postsongid = $2 `

        post.Likes = []responses.UserIdentifer{}
        post.Dislikes = []responses.UserIdentifer{}

        rows, err := executor.Query(query2, spotifyID, postID)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }

        for rows.Next() {
            userID := responses.UserIdentifer{}
            liked := true
            err := rows.Scan(&userID.SpotifyID, &userID.Username, &liked)
            if err != nil {
                return customerrors.WrapBasicError(err)
            }
            if liked {
                post.Likes = append(post.Likes, userID)
            } else {
                post.Dislikes = append(post.Dislikes, userID)
            }
        }

        return nil
}

func(p *PostsDAO) getUserPostPreviewsProperties(executor db.QueryExecutor, spotifyID string, createdAt time.Time, postPreviews []responses.PostPreview) error {
    query := `SELECT posts.albumarturi, posts.albumid, posts.albumname, posts.createdat, posts.rating, posts.songid, posts.songname, posts.review, posts.updatedat, posts.posterspotifyid, users.username
                FROM posts 
                INNER JOIN users 
                ON users.spotifyid = posts.posterspotifyid
                WHERE posts.posterspotifyid = $1 AND posts.createdat < $2 ORDER BY posts.createdat DESC LIMIT 25 `

        rows, err := executor.Query(query, spotifyID, createdAt)

        if err != nil {
            return customerrors.WrapBasicError(err)
        }


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
            postPreviews = append(postPreviews, post)
        }

        return nil
}

func(p *PostsDAO) getUserPostPreviews(executor db.QueryExecutor, spotifyID string, createdAt time.Time, paginationResponse *responses.PaginationResponse[[]responses.PostPreview, time.Time]) error {

    query := `SELECT spotifyid from USERS WHERE spotifyid = $1`

    res, err := executor.Exec(query, spotifyID)

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

    postPreviewsResponse := []responses.PostPreview{}
    err = p.getUserPostPreviewsProperties(executor, spotifyID, createdAt, postPreviewsResponse)

    if err != nil {
        return err
    }

    for i := 0; i < len(postPreviewsResponse); i++ {
        err = p.getPostLikesAndDislikes(executor, postPreviewsResponse[i].SpotifyID, postPreviewsResponse[i].SongID, &postPreviewsResponse[i])
        if err != nil {
            return err
        }
    }

    paginationResponse.DataResponse = postPreviewsResponse

    if len(postPreviewsResponse) > 0 {
        paginationResponse.PaginationKey = postPreviewsResponse[len(postPreviewsResponse)-1].CreatedAt
    } 

    return nil

}
