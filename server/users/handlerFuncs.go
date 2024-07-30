package users

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/gin-gonic/gin"
)

// GetUserByID godoc
// @Summary Gets a tunes user by their spotify ID
// @Description Gets a tunes user by their spotifyID
// @Tags users
// @Accept json
// @Produce json
// @Param spotifyID path string true "User Spotify ID"
// @Success 200 {object} responses.User 
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Router /users/{spotifyID} [get]
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
