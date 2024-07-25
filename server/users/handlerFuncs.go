package users

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/gin-gonic/gin"
)

func GetUserById(c *gin.Context) {

    spotifyID := c.Param("spotifyID")

    user, err := db.GetUserFromDbBySpotifyID(spotifyID)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, user) 
}

func UnFollowUser(c *gin.Context) {

    otherUserSpotifyID := c.Param("otherUserSpotifyID")
    spotifyID, found := c.Get("spotifyID")

    if otherUserSpotifyID == spotifyID {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Unfollowing not reflexive"})
        return
    }

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Middleware issue" })
        return
    }

    err := db.UnfollowUser(spotifyID.(string), otherUserSpotifyID)

    if err != nil {
        c.AbortWithError(-1, err)
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
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, followersPaginated)

}

func GetFollowers(c *gin.Context) {

    spotifyID, found := c.Get("spotifyID")
    paginationKey := c.Query("spotifyID")

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Jwt issue"})
        return
    }

    if paginationKey == "" {
        paginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
    }

    followersPaginated, err := db.GetFollowers(spotifyID.(string), paginationKey)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, followersPaginated)


}

func FollowerUser(c *gin.Context) {
    otherUserSpotifyID := c.Param("otherUserSpotifyID")
    spotifyID, found := c.Get("spotifyID")

    if otherUserSpotifyID == spotifyID {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Bad JSON body"})
        return
    }

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JWT fuckup"})
        return
    }

    err := db.FollowUser(spotifyID.(string), otherUserSpotifyID)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }


    c.Status(http.StatusNoContent)

}


func GetCurrentUser(c *gin.Context) {

    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !spotifyIdExists {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Jwtfuckup"})
        return
    }

    user, err := db.GetUserFromDbBySpotifyID(spotifyID.(string))

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, user)
}

func UpdateUserBySpotifyID(c *gin.Context) {

    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID := c.Param("spotifyID")

    if !found || spotifyID == "" {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JwtFuckup"})
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    if userUpdateRequest.Bio == nil && userUpdateRequest.Role == nil {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "give a body"})
        return
    }

    resp, e := updateUser(spotifyID, userUpdateRequest, userRole.(responses.Role))

    if e != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, resp)
}


func UpdateCurrentUserProperties(c *gin.Context) {
    
    userUpdateRequest := &requests.UpdateUserRequestDTO{}
    userRole, found := c.Get("userRole")
    spotifyID, spotifyIdExists := c.Get("spotifyID")

    if !found || !spotifyIdExists {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JwtFuckup"})
        return
    }

    err := c.ShouldBindBodyWithJSON(userUpdateRequest)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    if userUpdateRequest.Bio == nil && userUpdateRequest.Role == nil {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Bad body"})
        return
    }

    resp, e := updateUser(spotifyID.(string), userUpdateRequest, userRole.(responses.Role))

    if e != nil {
        c.AbortWithError(-1, e)
        return
    }


    c.JSON(http.StatusOK, resp)

}

func DeleteCurrentUser(c *gin.Context) {
    spotifyID, spotifyIdExists := c.Get("spotifyID")
    if !spotifyIdExists {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Bad"})
        return
    }

    err := db.DeleteUserByID(spotifyID.(string))

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.Status(http.StatusNoContent)
}

func DeleteUserBySpotifyID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")

    err := db.DeleteUserByID(spotifyID)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.Status(http.StatusNoContent)
}

func updateUser(spotifyID string, userUpdateRequest *requests.UpdateUserRequestDTO, userRole responses.Role) (*responses.User, error) {

    if userUpdateRequest.Role != nil && userRole != responses.ADMIN {
        return nil, customerrors.CustomError{StatusCode: http.StatusUnauthorized, Msg: "kill me last second last time"}
    }

    if userUpdateRequest.Role != nil && !responses.IsValidRole(*userUpdateRequest.Role) {
        return nil, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "get your shit together"}
    }

    return db.UpdateUserPropertiesBySpotifyID(spotifyID, userUpdateRequest)

}
