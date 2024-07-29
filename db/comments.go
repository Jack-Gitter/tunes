package db

import (
	"context"
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/responses"
)


func CreateComment(commentorID string, posterID string, songID string, commentText string) (*responses.Comment, error){

    query := `INSERT INTO comments (commentorspotifyid, posterspotifyid, songid, commenttext, likes, dislikes) values ($1, $2, $3, $4, $5, $6) RETURNING commentid, commentorspotifyid, posterspotifyid, songid, commenttext, likes, dislikes`

    res := DB.Driver.QueryRow(query, commentorID, posterID, songID, commentText, 0, 0)

    commentResp := &responses.Comment{}
    err := res.Scan(&commentResp.CommentID, &commentResp.CommentorID, &commentResp.PostSpotifyID, &commentResp.SongID, &commentResp.CommentText, &commentResp.Likes, &commentResp.Dislikes)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    return commentResp, nil


}

func DeleteComment(commentID string) error {
    query := `DELETE FROM comments WHERE commentid = $1`

    resp, err := DB.Driver.Exec(query, commentID)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    rows, err := resp.RowsAffected()

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    if rows < 1 {
        return &customerrors.CustomError{StatusCode: http.StatusNotFound, Msg: "resource not found"}
    }

    return nil

}

func DeleteCurrentUserComment(commentID string, spotifyID string) error {

    query := `DELETE FROM comments WHERE commentid = $1 AND commentorspotifyid = $2`

    resp, err := DB.Driver.Exec(query, commentID, spotifyID)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    rows, err := resp.RowsAffected()

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    if rows < 1 {
        return &customerrors.CustomError{StatusCode: http.StatusNotFound, Msg: "comment not found"}
    }

    return nil

}

func GetComment(commentID string) (*responses.Comment, error) {
    tx, err := DB.Driver.BeginTx(context.Background(), nil)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    defer tx.Rollback()

    _, err = tx.Exec(`SET TRANSACTION ISOLATION LEVEL REPEATABLE READ`)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    query := `SELECT comments.commentid, comments.commentorspotifyid, comments.posterspotifyid, comments.songid, comments.commenttext, users.username 
              FROM comments INNER JOIN users ON commentorspotifyid = spotifyid 
              WHERE commentid = $1`

    res := tx.QueryRow(query, commentID)

    commentResponse := &responses.Comment{}

    err = res.Scan(&commentResponse.CommentID, 
                &commentResponse.CommentorID, 
                &commentResponse.PostSpotifyID, 
                &commentResponse.SongID, 
                &commentResponse.CommentText,
                &commentResponse.CommentorUsername)
                

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    query = `SELECT liked FROM comment_votes WHERE commentid = $1`

    row, err := tx.Query(query, commentID)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    likes := 0
    dislikes := 0
    for row.Next() {
       val := true 
       err := row.Scan(&val)
       if err != nil {
           return nil, customerrors.WrapBasicError(err)
       }
       if val {
           likes +=1
       } else {
           dislikes+=1
       }
    }

    commentResponse.Likes = likes
    commentResponse.Dislikes = dislikes

    err = tx.Commit()

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    return commentResponse, nil

}

func LikeComment(commentID string, spotifyID string)  error {
    
    query := `INSERT INTO comment_votes (commentid, liked, voterspotifyid) values ($1, $2, $3) ON CONFLICT (commentid, voterspotifyid) DO UPDATE set liked = $2`

    res, err := DB.Driver.Exec(query, commentID, true, spotifyID)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    if rows < 1 {
        return customerrors.CustomError{StatusCode: http.StatusNotFound, Msg: "comment not found"}
    }

    return nil

}

func DislikeComment(commentID string, spotifyID string) error {
    
    query := `INSERT INTO comment_votes (commentid, liked, voterspotifyid) values ($1, $2, $3) ON CONFLICT (commentid, voterspotifyid) DO UPDATE SET liked = $2`

    res, err := DB.Driver.Exec(query, commentID, false, spotifyID)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    if rows < 1 {
        return customerrors.CustomError{StatusCode: http.StatusNotFound, Msg: "comment not found"}
    }

    return nil

}

func RemoveCommentVote(commentID string, spotifyID string) error {
    query := `DELETE FROM comment_votes WHERE commentid = $1 AND voterspotifyid = $2`

    res, err := DB.Driver.Exec(query, commentID, spotifyID)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    rows, err := res.RowsAffected()

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    if rows < 1 {
        return &customerrors.CustomError{StatusCode: http.StatusNotFound, Msg: "comment vote not found"}
    }

    return nil

}
