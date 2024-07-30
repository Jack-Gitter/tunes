package users

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Gets a tunes user by their spotify ID
// @Description Gets a tunes user by their spotifyID
// @Tags Users
// @Accept json
// @Produce json
// @Param spotifyID path string true "User Spotify ID"
// @Success 200 {object} responses.User 
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/{spotifyID} [get]
// @Security Bearer
func GetUserById(c *gin.Context) {

	spotifyID := c.Param("spotifyID")

	user, err := db.GetUserFromDbBySpotifyID(spotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Unfollowers a user for the currently signed in user
// @Description Unfollowers a user for the currently signed in user
// @Tags Users
// @Accept json
// @Produce json
// @Param spotifyID path string true "User spotify ID"
// @Param otherUserSpotifyID path string true "User to unfollow spotify ID"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/current/unfollow/{otherUserSpotifyID} [post]
// @Security Bearer
func UnFollowUser(c *gin.Context) {

	otherUserSpotifyID := c.Param("otherUserSpotifyID")
	spotifyID, found := c.Get("spotifyID")

	if otherUserSpotifyID == spotifyID {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Unfollowing not reflexive"})
		c.Abort()
		return
	}

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Middleware issue"})
		c.Abort()
		return
	}

	err := db.UnfollowUser(spotifyID.(string), otherUserSpotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Gets a users followers by their spotify ID
// @Description Gets a users followers by their spotify ID
// @Tags Users
// @Accept json
// @Produce json
// @Param spotifyID path string true "User spotify ID"
// @Param spotifyID query string false "Pagination Key for follow up responses. This key is a spotify ID"
// @Success 200 {object} responses.PaginationResponse[[]responses.User, string]
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/{spotifyID}/followers/ [get]
// @Security Bearer
func GetFollowersByID(c *gin.Context) {
	spotifyID := c.Param("spotifyID")
	paginationKey := c.Query("spotifyID")

	if paginationKey == "" {
		paginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
	}

	followersPaginated, err := db.GetFollowers(spotifyID, paginationKey)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, followersPaginated)

}

// @Summary Gets the current users followers
// @Description Gets the current users followers
// @Tags Users
// @Accept json
// @Produce json
// @Param spotifyID query string false "Pagination Key for follow up responses. This key is a spotify ID"
// @Success 200 {object} responses.PaginationResponse[[]responses.User, string]
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/current/followers/ [get]
// @Security Bearer
func GetFollowers(c *gin.Context) {

	spotifyID, found := c.Get("spotifyID")
	paginationKey := c.Query("spotifyID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Jwt issue"})
		c.Abort()
		return
	}

	if paginationKey == "" {
		paginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
	}

	followersPaginated, err := db.GetFollowers(spotifyID.(string), paginationKey)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, followersPaginated)

}

// @Summary Follows a user for the current user
// @Description Follows a user for the current user 
// @Tags Users
// @Accept json
// @Produce json
// @Param spotifyID path string false "Spotify ID of other user to follow"
// @Success 204 
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/current/follow/{otherUserSpotifyID} [post]
// @Security Bearer
func FollowerUser(c *gin.Context) {
	otherUserSpotifyID := c.Param("otherUserSpotifyID")
	spotifyID, found := c.Get("spotifyID")

	if otherUserSpotifyID == spotifyID {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Following is not reflexive"})
		c.Abort()
		return
	}

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JWT fuckup"})
		c.Abort()
		return
	}

	err := db.FollowUser(spotifyID.(string), otherUserSpotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Retreives the current user
// @Description Retrieves the current user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} responses.User
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/current [get]
// @Security Bearer
func GetCurrentUser(c *gin.Context) {

	spotifyID, spotifyIdExists := c.Get("spotifyID")

	if !spotifyIdExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Jwtfuckup"})
		c.Abort()
		return
	}

	user, err := db.GetUserFromDbBySpotifyID(spotifyID.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserBySpotifyID(c *gin.Context) {

	userUpdateRequest := &requests.UpdateUserRequestDTO{}
	spotifyID := c.Param("spotifyID")

	if spotifyID == "" {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JwtFuckup"})
		c.Abort()
		return
	}

	c.ShouldBindBodyWithJSON(userUpdateRequest)

	resp, err := updateUser(spotifyID, userUpdateRequest)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateCurrentUserProperties(c *gin.Context) {

	userUpdateRequest := &requests.UpdateUserRequestDTO{}
	spotifyID, spotifyIdExists := c.Get("spotifyID")

	if !spotifyIdExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JwtFuckup"})
		c.Abort()
		return
	}

	c.ShouldBindBodyWithJSON(userUpdateRequest)

	resp, e := updateUser(spotifyID.(string), userUpdateRequest)

	if e != nil {
		c.Error(e)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)

}

func DeleteCurrentUser(c *gin.Context) {
	spotifyID, spotifyIdExists := c.Get("spotifyID")
	if !spotifyIdExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Bad"})
		c.Abort()
		return
	}

	err := db.DeleteUserByID(spotifyID.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteUserBySpotifyID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")

	err := db.DeleteUserByID(spotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func updateUser(spotifyID string, userUpdateRequest *requests.UpdateUserRequestDTO) (*responses.User, error) {

	return db.UpdateUserPropertiesBySpotifyID(spotifyID, userUpdateRequest)

}
