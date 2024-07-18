package users

import (
	"fmt"
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/gin-gonic/gin"
)

func GetUserById(c *gin.Context) {

    spotifyID := c.Param("spotifyID")

    user, err := getUser(spotifyID)

    if err != nil {
        c.JSON(http.StatusBadRequest, err.Error())
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

    user, err := getUser(spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }

    c.JSON(http.StatusOK, user)
}

func UpdateUserBySpotifyID(c *gin.Context) {

    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID := c.Param("spotifyID")

    if !found || spotifyID == "" {
        c.JSON(http.StatusInternalServerError, "no role found for user")
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "invalid json body for updating a user!")
        return
    }

    if userUpdateRequest.Role != nil && userRole != responses.ADMIN {
        c.JSON(http.StatusUnauthorized, "cannot change your role if you're not admin!")
        return
    }

    if userUpdateRequest.Role != nil && !responses.IsValidRole(string(*userUpdateRequest.Role)) {
        c.JSON(http.StatusBadRequest, "invalid user role")
        return
    }

    user, found, err := db.UpdateUserPropertiesBySpotifyID(spotifyID, userUpdateRequest)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusBadRequest, "could not find user in db")
        return
    }

    c.JSON(http.StatusOK, user)
}

func UpdateCurrentUserProperties(c *gin.Context) {
    
    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !found || !spotifyIdExists {
        c.JSON(http.StatusInternalServerError, "no role found for user")
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "invalid json body for updating a user!")
        return
    }

    if userUpdateRequest.Role != nil && userRole != responses.ADMIN {
        c.JSON(http.StatusUnauthorized, "cannot change your role if you're not admin!")
        return
    }

    if userUpdateRequest.Role != nil && !responses.IsValidRole(string(*userUpdateRequest.Role)) {
        c.JSON(http.StatusBadRequest, "invalid user role")
        return
    }

    user, found, err := db.UpdateUserPropertiesBySpotifyID(spotifyID.(string), userUpdateRequest)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusBadRequest, "could not find user in db")
        return
    }

    c.JSON(http.StatusOK, user)

}

func getUser(spotifyID string) (*responses.User, error) {

    user, foundUser, err := db.GetUserFromDbBySpotifyID(spotifyID)

    if err != nil {
        return nil, err
    }

    if !foundUser {
        return nil, err
    }

    return user, nil

}
