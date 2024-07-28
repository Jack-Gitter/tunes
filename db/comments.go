package db

import (

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/responses"
)


func CreateComment(commentorID string, posterID string, songID string, commentText string) (*responses.Comment, error){

    query := `INSERT INTO comments (commentorspotifyid, posterspotifyid, songid, commenttext) values ($1, $2, $3, $4) RETURNING commentorspotifyid, posterspotifyid, songid, commenttext, likes, dislikes`

    res := DB.Driver.QueryRow(query, commentorID, posterID, songID, commentText)

    commentResp := &responses.Comment{}
    err := res.Scan(&commentResp.CommentorID, &commentResp.PostSpotifyID, &commentResp.SongID, &commentResp.CommentText, &commentResp.Likes, &commentResp.Dislikes)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    return commentResp, nil


}
