package users

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/gin-gonic/gin"
)

func GetUserById(c *gin.Context) {

    spotifyID := c.Query("spotifyID")

    if spotifyID == "" {
        c.JSON(http.StatusBadRequest, "please provide a user ID as a query parameter!")
        return
    }

    user, foundUser, err := db.GetUserFromDbBySpotifyID(spotifyID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !foundUser {
        c.JSON(http.StatusBadRequest, "no user with that ID has been found")
        return
    }

    c.JSON(http.StatusOK, user) 
}

func GetCurrentUser(c *gin.Context) {

    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !spotifyIdExists {
        c.JSON(http.StatusUnauthorized, "please log in again. no JWT found")
        return
    }

    user, foundUser, err := db.GetUserFromDbBySpotifyID(spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !foundUser {
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }

    c.JSON(http.StatusOK, user)
}
