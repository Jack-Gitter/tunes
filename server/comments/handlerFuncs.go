package comments

import (
	"net/http"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/gin-gonic/gin"
)

// @Summary Creates a comment for the current user
// @Description Creates a comment for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param CreatePostDTO body requests.CreateCommentDTO true "Information required to create a commment"
// @Param spotifyID path string true "spotifyID of poster"
// @Param songID path string true "songID of post to make a comment on"
// @Success 200 {object} responses.Comment
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/{spotifyID}/{songID} [post]
// @Security Bearer
func CreateComment(c *gin.Context) {

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

    comment, err := db.CreateComment(commentorID.(string), posterID, songID, createCommentDTO.CommentText)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, comment)

}

// @Summary Deletes a comment Must be admin
// @Description Deletes a comment. Must be admin
// @Tags Posts
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to delete"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/{commentID} [delete]
// @Security Bearer
func DeleteComment(c *gin.Context) {

    commentID := c.Param("commentID")

    err := db.DeleteComment(commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.Status(http.StatusNoContent)
}

// @Summary Deletes a comment for the current user
// @Description Deletes a comment for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to delete"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/current/{commentID} [delete]
// @Security Bearer
func DeleteCurrentUserComment(c *gin.Context) {

    commentID := c.Param("commentID")
    spotifyID, exists := c.Get("spotifyID")

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
        c.Abort()
        return
    }

    err := db.DeleteCurrentUserComment(commentID, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.Status(http.StatusNoContent)
}

// @Summary Retrieves a comment
// @Description Retrieves a comment
// @Tags Posts
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to retrieve"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/{commentID} [get]
// @Security Bearer
func GetComment(c *gin.Context)  {

    commentID := c.Param("commentID") 

    comment, err := db.GetComment(commentID)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, comment)

}

// @Summary Like a comment
// @Description Like a comment
// @Tags Posts
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to like"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/like/{commentID} [post]
// @Security Bearer
func LikeComment(c *gin.Context) {
    commentID := c.Param("commentID")
    spotifyID, exists := c.Get("spotifyID")

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuck"})
        c.Abort()
        return
    }


    err := db.LikeComment(commentID, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return 
    }

    c.Status(http.StatusNoContent)

}

// @Summary Dislike a comment
// @Description Dislike a comment
// @Tags Posts
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to dislike"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/dislike/{commentID} [post]
// @Security Bearer
func DislikeComment(c *gin.Context) {
    commentID := c.Param("commentID")
    spotifyID, exists := c.Get("spotifyID")

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuck"})
        c.Abort()
        return
    }


    err := db.DislikeComment(commentID, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return 
    }

    c.Status(http.StatusNoContent)
}

// @Summary Delete a vote on a comment for the current user
// @Description Delete a vote on a comment for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to remove the vote from"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/votes/current/{commentID} [delete]
// @Security Bearer
func RemoveCommentVote(c *gin.Context) {

    commentID := c.Param("commentID")
    spotifyID, exists := c.Get("spotifyID")

    if !exists {
        c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuck"})
        c.Abort()
        return
    }

    err := db.RemoveCommentVote(commentID, spotifyID.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }
    
    c.Status(http.StatusNoContent)

}

// @Summary Updates a comment for the current user
// @Description Updates a comment for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param commentID path string true "Comment ID of comment to update"
// @Param UpdateCommentDTO body requests.UpdateCommentDTO true "Comment data to update"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /comments/current/{commentID} [patch]
// @Security Bearer
func UpdateComment(c *gin.Context) {

    commentID := c.Param("commentID")
    updateCommentDTO := &requests.UpdateCommentDTO{}
    c.ShouldBindBodyWithJSON(updateCommentDTO)
    

    resp, err := db.UpdateComment(commentID, updateCommentDTO)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, resp)

}
