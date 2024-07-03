package posts

import (
	//"github.com/Jack-Gitter/tunes/db"
	"net/http"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/gin-gonic/gin"
)


func CreatePost(c *gin.Context) {

    post := &models.Post{}
    err := c.ShouldBindBodyWithJSON(post)
    if err != nil {
        c.JSON(http.StatusBadRequest, "invalid post request body for request 'post'")
    }
    // make the middleware to grab the user spotify id from the jwt
    //db.CreatePost(post)
     
}
