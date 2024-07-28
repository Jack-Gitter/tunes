package db

import (
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
