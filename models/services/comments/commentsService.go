package comments

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/daos"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/gin-gonic/gin"
)
type CommentsService struct {
    DB *sql.DB
    CommentsDAO daos.ICommentsDAO
}

type ICommentsService interface {
    CreateComment(c *gin.Context) 
    DeleteComment(c *gin.Context) 
    DeleteCurrentUserComment(c *gin.Context) 
    GetComment(c *gin.Context)  
    LikeComment(c *gin.Context) 
    DislikeComment(c *gin.Context) 
    RemoveCommentVote(c *gin.Context) 
    UpdateComment(c *gin.Context) 
}

// @Summary Creates a comment for the current user
// @Description Creates a comment for the current user
// @Tags Comments
// @Accept json
// @Produce json
// @Param CreatePostDTO body requests.CreateCommentDTO true "Information required to create a commment"
// @Param spotifyID path string true "spotifyID of poster"
// @Param songID path string true "songID of post to make a comment on"
// @Success 200 {object} responses.Comment
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/{spotifyID}/{songID} [post]
// @Security Bearer
func(cs *CommentsService) CreateComment(c *gin.Context) {

    commentorID, exists := c.Get("spotifyID")
    posterID := c.Param("spotifyID")
    songID := c.Param("songID")

    createCommentDTO := &requests.CreateCommentDTO{}

    c.ShouldBindBodyWithJSON(createCommentDTO)

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusBadGateway, Msg: "JWT elmsss"})
        c.Abort()
        return
    }

    comment, err := cs.CommentsDAO.CreateComment(cs.DB, commentorID.(string), posterID, songID, createCommentDTO.CommentText)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, comment)

}

// @Summary Deletes a comment Must be admin
// @Description Deletes a comment. Must be admin
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to delete"
// @Success 204
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 403 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/admin/{commentID} [delete]
// @Security Bearer
func(cs *CommentsService) DeleteComment(c *gin.Context) {

    commentID := c.Param("commentID")

    err := cs.CommentsDAO.DeleteComment(cs.DB, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.Status(http.StatusNoContent)
}

// @Summary Deletes a comment for the current user
// @Description Deletes a comment for the current user
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to delete"
// @Success 204
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/current/{commentID} [delete]
// @Security Bearer
func(cs *CommentsService) DeleteCurrentUserComment(c *gin.Context) {

    commentID := c.Param("commentID")

    // need some extra validation here, first need to get the comment. then see if the spotifyid for the comment is the same as our spotifyid
    err := cs.CommentsDAO.DeleteComment(cs.DB, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.Status(http.StatusNoContent)
}

// @Summary Retrieves a comment
// @Description Retrieves a comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to retrieve"
// @Success 204
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/{commentID} [get]
// @Security Bearer
func(cs *CommentsService) GetComment(c *gin.Context)  {

    commentID := c.Param("commentID") 

    tx, err := cs.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    comment, err := cs.CommentsDAO.GetCommentProperties(tx, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    likes, dislikes, err := cs.CommentsDAO.GetCommentVotes(tx, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    comment.Likes = len(likes)
    comment.Dislikes = len(dislikes)

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }


    c.JSON(http.StatusOK, comment)

}

// @Summary Like a comment
// @Description Like a comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to like"
// @Success 204
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 409 {string} string 
// @Failure 500 {string} string 
// @Router /comments/like/{commentID} [post]
// @Security Bearer
func(cs *CommentsService) LikeComment(c *gin.Context) {
    commentID := c.Param("commentID")
    spotifyID, exists := c.Get("spotifyID")

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuck"})
        c.Abort()
        return
    }

    tx, err := cs.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    likes, _, err := cs.CommentsDAO.GetCommentVotes(tx, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    for _, userIdentifier := range likes {
        if userIdentifier.SpotifyID == spotifyID {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusConflict, Msg: "cannot like a message twice"})
            c.Abort()
            return
        }
    }

    err = cs.CommentsDAO.LikeComment(tx, commentID, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    c.Status(http.StatusNoContent)

}

// @Summary Dislike a comment
// @Description Dislike a comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to dislike"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 409 {string} string 
// @Failure 500 {string} string 
// @Router /comments/dislike/{commentID} [post]
// @Security Bearer
func(cs *CommentsService) DislikeComment(c *gin.Context) {
    commentID := c.Param("commentID")
    spotifyID, exists := c.Get("spotifyID")

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuck"})
        c.Abort()
        return
    }

    tx, err := cs.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    _, dislikes, err := cs.CommentsDAO.GetCommentVotes(tx, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    for _, userIdentifier := range dislikes {
        if userIdentifier.SpotifyID == spotifyID {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusConflict, Msg: "cannot like a message twice"})
            c.Abort()
            return
        }
    }

    err = cs.CommentsDAO.DislikeComment(tx, commentID, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    c.Status(http.StatusNoContent)
}

// @Summary Delete a vote on a comment for the current user
// @Description Delete a vote on a comment for the current user
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to remove the vote from"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/votes/current/{commentID} [delete]
// @Security Bearer
func(cs *CommentsService) RemoveCommentVote(c *gin.Context) {

    commentID := c.Param("commentID")
    spotifyID, exists := c.Get("spotifyID")

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuck"})
        c.Abort()
        return
    }

    err := cs.CommentsDAO.RemoveCommentVote(cs.DB, commentID, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }
    
    c.Status(http.StatusNoContent)

}

// @Summary Updates a comment for the current user
// @Description Updates a comment for the current user
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to update"
// @Param UpdateCommentDTO body requests.UpdateCommentDTO true "Comment data to update"
// @Success 204
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/current/{commentID} [patch]
// @Security Bearer
func(cs *CommentsService) UpdateComment(c *gin.Context) {

    commentID := c.Param("commentID")
    updateCommentDTO := &requests.UpdateCommentDTO{}
    c.ShouldBindBodyWithJSON(updateCommentDTO)

    tx, err := cs.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

    err = db.SetTransactionIsolationLevel(tx, sql.LevelRepeatableRead)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    resp, err := cs.CommentsDAO.UpdateComment(tx, commentID, updateCommentDTO)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    likes, dislikes, err := cs.CommentsDAO.GetCommentVotes(tx, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    newcomment, err := cs.CommentsDAO.GetCommentProperties(tx, commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    resp.CommentorUsername = newcomment.CommentorUsername
    resp.Likes = len(likes)
    resp.Dislikes = len(dislikes)

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, resp)

}
