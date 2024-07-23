package users

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/gin-gonic/gin"
)

func GetUserById(c *gin.Context) {

    spotifyID := c.Param("spotifyID")

    user, found, err := db.GetUserFromDbBySpotifyID(spotifyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "No user with specified ID found")
        return
    }

    c.JSON(http.StatusOK, user) 
}

func UnFollowUser(c *gin.Context) {

    otherUserSpotifyID := c.Param("otherUserSpotifyID")
    spotifyID, found := c.Get("spotifyID")

    if otherUserSpotifyID == spotifyID {
        c.JSON(http.StatusBadRequest, "Unfollowing is not reflexive")
        return
    }

    if !found {
        c.JSON(http.StatusInternalServerError, "spotifyID key not set from JWT middleware")
        return
    }

    found, err := db.UnfollowUser(spotifyID.(string), otherUserSpotifyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "One of the provided spotifyIDs failed to map to a valid user")
        return
    }

    c.Status(http.StatusNoContent)

}

func GetFollowersByID(c *gin.Context) {
    spotifyID := c.Param("spotifyID")
    paginationKey := c.Query("spotifyID")

    if paginationKey == "" {
        paginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
    }

    followersPaginated, found, err := db.GetFollowers(spotifyID, paginationKey)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "The provided spotifyID failed to map to a valid user")
        return
    }

    c.JSON(http.StatusOK, followersPaginated)

}

func GetFollowers(c *gin.Context) {

    spotifyID, found := c.Get("spotifyID")
    paginationKey := c.Query("spotifyID")

    if !found {
        c.JSON(http.StatusInternalServerError, "No spotifyID variable set from JWT middleware")
        return
    }

    if paginationKey == "" {
        paginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
    }

    followersPaginated, found, err := db.GetFollowers(spotifyID.(string), paginationKey)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "The provided spotifyID failed to map to a valid user")
        return
    }

    c.JSON(http.StatusOK, followersPaginated)


}

func FollowerUser(c *gin.Context) {
    otherUserSpotifyID := c.Param("otherUserSpotifyID")
    spotifyID, found := c.Get("spotifyID")

    if otherUserSpotifyID == spotifyID {
        c.JSON(http.StatusBadRequest, "Following is not reflexive")
        return
    }

    if !found {
        c.JSON(http.StatusInternalServerError, "No spotifyID set from JWT middleware")
        return
    }

    found, err := db.FollowUser(spotifyID.(string), otherUserSpotifyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "One of the provided spotifyIDs failed to map to a valid user")
        return
    }

    c.Status(http.StatusNoContent)

}


func GetCurrentUser(c *gin.Context) {

    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !spotifyIdExists {
        c.JSON(http.StatusInternalServerError, "No spotifyID key set from JWT middleware")
        return
    }

    user, found, err := db.GetUserFromDbBySpotifyID(spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "No user with ID found")
        return
    }

    c.JSON(http.StatusOK, user)
}

func UpdateUserBySpotifyID(c *gin.Context) {

    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID := c.Param("spotifyID")

    if !found || spotifyID == "" {
        c.JSON(http.StatusInternalServerError, "No role variable set by JWT middleware")
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "Invalid JSON body")
        return
    }

    resp, found, err := updateUser(spotifyID, userUpdateRequest, userRole.(responses.Role))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "No user found with provided spotifyID")
        return
    }

    c.JSON(http.StatusOK, resp)
}


func UpdateCurrentUserProperties(c *gin.Context) {
    
    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !found || !spotifyIdExists {
        c.JSON(http.StatusInternalServerError, "No role variable set by JWT middleware")
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "Invalid JSON body")
        return
    }

    resp, found, err := updateUser(spotifyID.(string), userUpdateRequest, userRole.(responses.Role))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "No user with provided spotifyID found")
        return
    }

    c.JSON(http.StatusOK, resp)

}

func DeleteCurrentUser(c *gin.Context) {
    spotifyID, spotifyIdExists := c.Get("spotifyID")
    if !spotifyIdExists {
        c.JSON(http.StatusInternalServerError, "No spotifyID variable set by JWT middleware")
        return
    }

    found, err := db.DeleteUserByID(spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "No user found with provided spotifyID")
        return
    }

    c.Status(http.StatusNoContent)
}

func DeleteUserBySpotifyID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")

    found, err := db.DeleteUserByID(spotifyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "No user found with provided spotifyID")
        return
    }

    c.Status(http.StatusNoContent)
}

func updateUser(spotifyID string, userUpdateRequest *requests.UpdateUserRequestDTO, userRole responses.Role) (*responses.User, bool, error) {

    if userUpdateRequest.Role != nil && userRole != responses.ADMIN {
        return nil, false, errors.New("Do not have sufficient permissions to change roles")
    }

    if userUpdateRequest.Role != nil && !responses.IsValidRole(*userUpdateRequest.Role) {
        return nil, false, errors.New("User rle provided is not valid")
    }

    resp, found, err := db.UpdateUserPropertiesBySpotifyID(spotifyID, userUpdateRequest)

    if err != nil {
        return nil, false, err
    }

    if !found {
        return nil, false, nil
    }

    return resp, true, nil

}
