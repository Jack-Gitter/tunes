package users

import (
	"net/http"

	"github.com/Jack-Gitter/tunes/customerrors"
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
        if tunesError, ok := err.(customerrors.TunesError); ok {
            if tunesError.ErrorType == customerrors.NoDatabaseRecordsFoundError {
                c.JSON(http.StatusBadRequest, err.Error())
                return
            } else if tunesError.ErrorType == customerrors.Neo4jDatabaseRequestError {
                c.JSON(http.StatusInternalServerError, err.Error())
                return
            }
        } 
    }

    c.JSON(http.StatusOK, user) 
}

func GetCurrentUser(c *gin.Context) {

    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !spotifyIdExists {
        c.JSON(http.StatusUnauthorized, "please log in again. no JWT found")
        return
    }

    user, err := db.GetUserFromDbBySpotifyID(spotifyID.(string))

    if err != nil {
        if tunesError, ok := err.(customerrors.TunesError); ok {
            if tunesError.ErrorType == customerrors.NoDatabaseRecordsFoundError {
                c.JSON(http.StatusBadRequest, err.Error())
                return
            } else if tunesError.ErrorType == customerrors.Neo4jDatabaseRequestError {
                c.JSON(http.StatusInternalServerError, err.Error())
                return
            }
        } 
    }

    c.JSON(http.StatusOK, user)
}
