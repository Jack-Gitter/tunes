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

    user, err := db.GetUserFromDbBySpotifyID(spotifyID)

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, user) 
}

func UnFollowUser(c *gin.Context) {

    otherUserSpotifyID := c.Param("otherUserSpotifyID")
    spotifyID, found := c.Get("spotifyID")

    if otherUserSpotifyID == spotifyID {
        /*c.AbortWithError(http.StatusInternalServerError, badrequest)
        return*/
        c.JSON(http.StatusBadRequest, "Unfollowing is not reflexive")
        return
    }

    if !found {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "spotifyID key not set from JWT middleware")
        return
    }

    err := db.UnfollowUser(spotifyID.(string), otherUserSpotifyID)

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
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

    followersPaginated, err := db.GetFollowers(spotifyID, paginationKey)

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }


    c.JSON(http.StatusOK, followersPaginated)

}

func GetFollowers(c *gin.Context) {

    spotifyID, found := c.Get("spotifyID")
    paginationKey := c.Query("spotifyID")

    if !found {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "No spotifyID variable set from JWT middleware")
        return
    }

    if paginationKey == "" {
        paginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
    }

    followersPaginated, err := db.GetFollowers(spotifyID.(string), paginationKey)

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    if !found {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusNotFound, "The provided spotifyID failed to map to a valid user")
        return
    }

    c.JSON(http.StatusOK, followersPaginated)


}

func FollowerUser(c *gin.Context) {
    otherUserSpotifyID := c.Param("otherUserSpotifyID")
    spotifyID, found := c.Get("spotifyID")

    if otherUserSpotifyID == spotifyID {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusBadRequest, "Following is not reflexive")
        return
    }

    if !found {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "No spotifyID set from JWT middleware")
        return
    }

    err := db.FollowUser(spotifyID.(string), otherUserSpotifyID)

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }


    c.Status(http.StatusNoContent)

}


func GetCurrentUser(c *gin.Context) {

    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !spotifyIdExists {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "No spotifyID key set from JWT middleware")
        return
    }

    user, err := db.GetUserFromDbBySpotifyID(spotifyID.(string))

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }


    c.JSON(http.StatusOK, user)
}

func UpdateUserBySpotifyID(c *gin.Context) {

    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID := c.Param("spotifyID")

    if !found || spotifyID == "" {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "No role variable set by JWT middleware")
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "Invalid JSON body")
        return
    }

    if userUpdateRequest.Bio == nil && userUpdateRequest.Role == nil {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusBadRequest, "Must provide at least one parameter to chage")
        return
    }

    resp, e := updateUser(spotifyID, userUpdateRequest, userRole.(responses.Role))

    if e != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }


    c.JSON(http.StatusOK, resp)
}


func UpdateCurrentUserProperties(c *gin.Context) {
    
    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !found || !spotifyIdExists {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "No role variable set by JWT middleware")
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "Invalid JSON body")
        return
    }

    if userUpdateRequest.Bio == nil && userUpdateRequest.Role == nil {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusBadRequest, "Must provide at least one parameter to chage")
        return
    }

    resp, e := updateUser(spotifyID.(string), userUpdateRequest, userRole.(responses.Role))

    if e != nil {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        if err, ok := err.(*db.DBError); ok {
            c.JSON(err.StatusCode, err.Msg)
        } else {
            c.JSON(http.StatusInternalServerError, err.Error())
        }
        return
    }


    c.JSON(http.StatusOK, resp)

}

func DeleteCurrentUser(c *gin.Context) {
    spotifyID, spotifyIdExists := c.Get("spotifyID")
    if !spotifyIdExists {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "No spotifyID variable set by JWT middleware")
        return
    }

    err := db.DeleteUserByID(spotifyID.(string))

    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.Status(http.StatusNoContent)
}

func DeleteUserBySpotifyID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")

    err := db.DeleteUserByID(spotifyID)

    if err != nil {
        /*c.AbortWithError(http.StatusInternalServerError, internalserver)
        return*/
        c.JSON(http.StatusInternalServerError, "No spotifyID variable set by JWT middleware")
        return
    }

    c.Status(http.StatusNoContent)
}

func updateUser(spotifyID string, userUpdateRequest *requests.UpdateUserRequestDTO, userRole responses.Role) (*responses.User, *db.DBError) {

    if userUpdateRequest.Role != nil && userRole != responses.ADMIN {
        return nil, &db.DBError{StatusCode: http.StatusUnauthorized, Msg: "Do not have sufficient permissions to change roles"}
    }

    if userUpdateRequest.Role != nil && !responses.IsValidRole(*userUpdateRequest.Role) {
        return nil, &db.DBError{StatusCode: http.StatusBadRequest, Msg: "role provided is not valid"}
    }

    return db.UpdateUserPropertiesBySpotifyID(spotifyID, userUpdateRequest)


}
