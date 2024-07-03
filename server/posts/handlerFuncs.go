package posts

import (
	"net/http"

	//"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/gin-gonic/gin"
)


func CreatePostForCurrentUser(c *gin.Context) {

    tunesPost := &models.Post{}

    err := c.ShouldBindBodyWithJSON(tunesPost)

    if err != nil {
        c.JSON(http.StatusBadRequest, "invalid post request body for request 'post'")
    }

    // make the middleware to grab the user spotify id from the jwt
   // db.CreatePost(*tunesPost, )
     
}
