package daos

import (
	"net/http"
	"time"
	"github.com/Jack-Gitter/tunes/db"
	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/mitchellh/mapstructure"
)

type CommentsDAO struct { }

type ICommentsDAO interface {
    CreateComment(executor db.QueryExecutor, commentorID string, posterID string, songID string, commentText string) (*responses.Comment, error)
    DeleteComment(executor db.QueryExecutor, commentID string) error
    DeleteCurrentUserComment(executor db.QueryExecutor, commentID string, spotifyID string) error 
    GetComment(executor db.QueryExecutor, commentID string) (*responses.Comment, error) 
    LikeOrDislikeComment(executor db.QueryExecutor, commentID string, spotifyID string, liked bool)  error 
    RemoveCommentVote(executor db.QueryExecutor, commentID string, spotifyID string) error 
    UpdateComment(executor db.QueryExecutor, commentID string, updateCommentDTO *requests.UpdateCommentDTO) (*responses.Comment, error) 
}

func(c *CommentsDAO) CreateComment(executor db.QueryExecutor, commentorID string, posterID string, songID string, commentText string) (*responses.Comment, error){

    query := `INSERT INTO comments (commentorspotifyid, posterspotifyid, songid, commenttext, createdAt, updatedAt) values ($1, $2, $3, $4, $5, $5) RETURNING commentid, commentorspotifyid, posterspotifyid, songid, commenttext`

    res := executor.QueryRow(query, commentorID, posterID, songID, commentText, time.Now().UTC())

    commentResp := &responses.Comment{}
    err := res.Scan(&commentResp.CommentID, &commentResp.CommentorID, &commentResp.PostSpotifyID, &commentResp.SongID, &commentResp.CommentText)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    return commentResp, nil


}

func(c *CommentsDAO) DeleteComment(executor db.QueryExecutor, commentID string) error {
    query := `DELETE FROM comments WHERE commentid = $1`

    resp, err := executor.Exec(query, commentID)

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

func(c *CommentsDAO) DeleteCurrentUserComment(executor db.QueryExecutor, commentID string, spotifyID string) error {

    query := `DELETE FROM comments WHERE commentid = $1 AND commentorspotifyid = $2`

    resp, err := executor.Exec(query, commentID, spotifyID)

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

func(c *CommentsDAO) GetComment(executor db.QueryExecutor, commentID string) (*responses.Comment, error) {

    commentResponse := &responses.Comment{}

    query := `SELECT comments.commentid, comments.commentorspotifyid, comments.posterspotifyid, comments.songid, comments.commenttext, comments.createdat, comments.updatedat, users.username 
              FROM comments INNER JOIN users ON commentorspotifyid = spotifyid 
              WHERE commentid = $1`

    res := executor.QueryRow(query, commentID)

    err := res.Scan(&commentResponse.CommentID, 
                &commentResponse.CommentorID, 
                &commentResponse.PostSpotifyID, 
                &commentResponse.SongID, 
                &commentResponse.CommentText,
                &commentResponse.CreatedAt,
                &commentResponse.UpdatedAt,
                &commentResponse.CommentorUsername)
                

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    query = `SELECT liked FROM comment_votes WHERE commentid = $1`

    row, err := executor.Query(query, commentID)

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

    return commentResponse, nil

}

func(c *CommentsDAO) LikeOrDislikeComment(executor db.QueryExecutor, commentID string, spotifyID string, liked bool)  error {
    

    query := `SELECT COUNT(*) FROM comment_votes WHERE commentid = $1 AND liked = $2 AND voterspotifyid = $3`

    row := executor.QueryRow(query, commentID, liked, spotifyID)

    count := 0
    err := row.Scan(&count)

    if err != nil {
        return customerrors.WrapBasicError(err)
    }

    if count >= 1 {
        return &customerrors.CustomError{StatusCode: http.StatusConflict, Msg: "cannot like/dislike comment twice"}
    }

    query = `INSERT INTO comment_votes (commentid, liked, voterspotifyid) values ($1, $2, $3) ON CONFLICT (commentid, voterspotifyid) DO UPDATE set liked = $2`

    res, err := executor.Exec(query, commentID, liked, spotifyID)

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


func(c *CommentsDAO) RemoveCommentVote(executor db.QueryExecutor, commentID string, spotifyID string) error {
    query := `DELETE FROM comment_votes WHERE commentid = $1 AND voterspotifyid = $2`

    res, err := executor.Exec(query, commentID, spotifyID)

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

func(c *CommentsDAO) UpdateComment(executor db.QueryExecutor, commentID string, updateCommentDTO *requests.UpdateCommentDTO) (*responses.Comment, error) {

    updateCommentMap := make(map[string]any)
    mapstructure.Decode(updateCommentDTO, &updateCommentMap)

    t := time.Now().UTC()
    updateCommentMap["updatedat"] = &t
    

    conditionals := make(map[string]any)
    conditionals["commentid"] = commentID

    returning := []string{"commentid", "commentorspotifyid", "posterspotifyid", "songid", "commenttext", "createdat", "updatedat"}

    query, vals := db.PatchQueryBuilder("comments", updateCommentMap, conditionals, returning)

    comment := &responses.Comment{}


    row := executor.QueryRow(query, vals...)

    err := row.Scan(&comment.CommentID, &comment.CommentorID, &comment.PostSpotifyID, &comment.SongID, &comment.CommentText, &comment.CreatedAt, &comment.UpdatedAt)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    query = `SELECT liked FROM comment_votes WHERE commentid = $1`

    rows, err := executor.Query(query, commentID)

    if err != nil {
        return nil, customerrors.WrapBasicError(err)
    }

    likes := 0
    dislikes := 0
    for rows.Next() {
       val := true 
       err := rows.Scan(&val)
       if err != nil {
           return nil, customerrors.WrapBasicError(err)
       }
       if val {
           likes +=1
       } else {
           dislikes+=1
       }
    }
    comment.Likes = likes
    comment.Dislikes = dislikes


    return comment, nil

}
