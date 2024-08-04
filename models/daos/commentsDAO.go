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
    GetCommentProperties(executor db.QueryExecutor, commentID string) (*responses.Comment, error) 
    GetCommentLikes(executor db.QueryExecutor, commentID string) (int, int, error)
    LikeComment(executor db.QueryExecutor, commentID string, spotifyID string) error 
    DislikeComment(executor db.QueryExecutor, commentID string, spotifyID string) error 
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

func(c *CommentsDAO) GetCommentProperties(executor db.QueryExecutor, commentID string) (*responses.Comment, error) {

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


    return commentResponse, nil

}

func(cs *CommentsDAO) GetCommentLikes(executor db.QueryExecutor, commentID string) (int, int, error) {
    query := `SELECT liked FROM comment_votes WHERE commentid = $1`

    row, err := executor.Query(query, commentID)

    if err != nil {
        return 0, 0, customerrors.WrapBasicError(err)
    }

    likes := 0
    dislikes := 0
    for row.Next() {
       val := true 
       err := row.Scan(&val)
       if err != nil {
           return 0, 0, customerrors.WrapBasicError(err)
       }
       if val {
           likes +=1
       } else {
           dislikes+=1
       }
    }

    return likes, dislikes, nil

}

func(c *CommentsDAO) LikeComment(executor db.QueryExecutor, commentID string, spotifyID string) error {
    
    query := `INSERT INTO comment_votes (commentid, liked, voterspotifyid) values ($1, $2, $3) ON CONFLICT (commentid, voterspotifyid) DO UPDATE set liked = $2`

    res, err := executor.Exec(query, commentID, true, spotifyID)

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

func(c *CommentsDAO) DislikeComment(executor db.QueryExecutor, commentID string, spotifyID string) error {

    query := `INSERT INTO comment_votes (commentid, liked, voterspotifyid) values ($1, $2, $3) ON CONFLICT (commentid, voterspotifyid) DO UPDATE SET liked = $2`

    res, err := executor.Exec(query, commentID, false, spotifyID)

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

    return comment, nil

}
