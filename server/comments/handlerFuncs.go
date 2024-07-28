package comments

import (
	"net/http"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/gin-gonic/gin"
)

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
