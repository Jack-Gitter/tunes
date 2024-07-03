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
