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

    user, err := db.GetUserFromDbBySpotifyID(spotifyID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to fetch the user from the database fml")
        return
    }

    c.JSON(http.StatusOK, user) 
}

func GetCurrentUser(c *gin.Context) {

    spotifyID, e1 := c.Get("spotifyID")

    if !e1 {
        c.JSON(http.StatusUnauthorized, "please log in again. no JWT found")
        return
    }

    user, err := db.GetUserFromDbBySpotifyID(spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to fetch the current user from the database")
        return
    }

    c.JSON(http.StatusOK, user)
}
