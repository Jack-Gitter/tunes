package posts

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/gin-gonic/gin"
)


func CreatePostForCurrentUser(c *gin.Context) {

    spotifyID, found := c.Get("spotifyID")

    if !found {
        c.JSON(http.StatusInternalServerError, "forgot to put the spotify id to the middleware")
        return
    }

    tunesPost := &models.Post{}

    err := c.ShouldBindBodyWithJSON(tunesPost)

    if err != nil {
        c.JSON(http.StatusBadRequest, "invalid post request body for request 'post'")
        return
    }

    // make the middleware to grab the user spotify id from the jwt
    err = db.CreatePost(tunesPost, spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, "failed to put post into db for some reason")
        return
    }

    c.JSON(http.StatusOK, tunesPost)

}
