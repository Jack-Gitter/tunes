package users

import (
	"context"
	"database/sql"
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/daos"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/gin-gonic/gin"
)

type UserService struct {
    DB *sql.DB
    UsersDAO daos.IUsersDAO
}

type IUserSerivce interface {
    GetUserById(c *gin.Context)
    UnFollowUser(c *gin.Context)
    GetFollowersByID(c *gin.Context)
    GetFollowers(c *gin.Context) 
    FollowUser(c *gin.Context)
    GetCurrentUser(c *gin.Context) 
    UpdateUserBySpotifyID(c *gin.Context)
    UpdateCurrentUserProperties(c *gin.Context) 
    DeleteCurrentUser(c *gin.Context)
    DeleteUserBySpotifyID(c *gin.Context)
}

// @Summary Gets a tunes user by their spotify ID
// @Description Gets a tunes user by their spotifyID
// @Tags Users
// @Accept json
// @Produce json
// @Param spotifyID path string true "User Spotify ID"
// @Success 200 {object} responses.User
// @Failure 401 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users/{spotifyID} [get]
// @Security Bearer
func(u *UserService) GetUserById(c *gin.Context) {

	spotifyID := c.Param("spotifyID")

	user, err := u.UsersDAO.GetUserFromDbBySpotifyID(u.DB, spotifyID)

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
// @Param otherUserSpotifyID path string true "User to unfollow spotify ID"
// @Success 204
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/current/unfollow/{otherUserSpotifyID} [delete]
// @Security Bearer
func(u *UserService) UnFollowUser(c *gin.Context) {

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

	err := u.UsersDAO.UnfollowUser(u.DB, spotifyID.(string), otherUserSpotifyID)

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
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/{spotifyID}/followers/ [get]
// @Security Bearer
func(u *UserService) GetFollowersByID(c *gin.Context) {
	spotifyID := c.Param("spotifyID")
	paginationKey := c.Query("spotifyID")

	if paginationKey == "" {
		paginationKey = "zzzzzzzzzzzzzzzzzzzzzzzzzz"
	}

	followersPaginated, err := u.UsersDAO.GetFollowers(spotifyID, paginationKey)

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
// @Failure 401 {string} string 
// @Failure 500 {string} string 
// @Router /users/current/followers/ [get]
// @Security Bearer
func(u *UserService) GetFollowers(c *gin.Context) {

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

    tx, err := u.DB.BeginTx(context.Background(), nil)

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
        c.Abort()
        return
    }

    defer tx.Rollback()

	followersPaginated, err := u.UsersDAO.GetFollowers(tx, spotifyID.(string), paginationKey)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

    err = tx.Commit()

    if err != nil {
        c.Error(customerrors.WrapBasicError(err))
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
// @Param otherUserSpotifyID path string true "Spotify ID of other user to follow"
// @Success 204 
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 409 {string} string 
// @Failure 500 {string} string 
// @Router /users/current/follow/{otherUserSpotifyID} [post]
// @Security Bearer
func(u *UserService) FollowUser(c *gin.Context) {
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

	err := u.UsersDAO.FollowUser(u.DB, spotifyID.(string), otherUserSpotifyID)

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
// @Failure 401 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/current [get]
// @Security Bearer
func(u *UserService) GetCurrentUser(c *gin.Context) {

	spotifyID, spotifyIdExists := c.Get("spotifyID")

	if !spotifyIdExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Jwtfuckup"})
		c.Abort()
		return
	}

	user, err := u.UsersDAO.GetUserFromDbBySpotifyID(u.DB, spotifyID.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Updates a user by their spotify ID. Only accessable to admins
// @Description Updates a user by their spotify ID. Only accessable to admins
// @Tags Users
// @Accept json
// @Produce json
// @Param UpdateUserDTO body requests.UpdateUserRequestDTO true "Information to update"
// @Param spotifyID path string true "Spotify ID of the user to update"
// @Success 200 {object} responses.User
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 403 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/admin/{spotifyID} [patch]
// @Security Bearer
func(u *UserService) UpdateUserBySpotifyID(c *gin.Context) {

	userUpdateRequest := &requests.UpdateUserRequestDTO{}
	spotifyID := c.Param("spotifyID")

	if spotifyID == "" {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JwtFuckup"})
		c.Abort()
		return
	}

	c.ShouldBindBodyWithJSON(userUpdateRequest)

	resp, err := u.updateUser(spotifyID, userUpdateRequest)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}


// @Summary Updates the current user
// @Description Updates the current user
// @Tags Users
// @Accept json
// @Produce json
// @Param UpdateUserDTO body requests.UpdateUserRequestDTO true "Information to update"
// @Success 200 {object} responses.User
// @Failure 400 {string} string 
// @Failure 401 {string} string 
// @Failure 403 {string} string 
// @Failure 500 {string} string 
// @Router /users/current [patch]
// @Security Bearer
func(u *UserService) UpdateCurrentUserProperties(c *gin.Context) {

	userUpdateRequest := &requests.UpdateUserRequestDTO{}
	spotifyID, spotifyIdExists := c.Get("spotifyID")

	if !spotifyIdExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "JwtFuckup"})
		c.Abort()
		return
	}

	c.ShouldBindBodyWithJSON(userUpdateRequest)

	resp, e := u.updateUser(spotifyID.(string), userUpdateRequest)

	if e != nil {
		c.Error(e)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)

}


// @Summary Deletes the current user
// @Description Deletes the current user
// @Tags Users
// @Accept json
// @Produce json
// @Success 204
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/current [delete]
// @Security Bearer
func(u *UserService) DeleteCurrentUser(c *gin.Context) {
	spotifyID, spotifyIdExists := c.Get("spotifyID")
	if !spotifyIdExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Bad"})
		c.Abort()
		return
	}

	err := u.UsersDAO.DeleteUserByID(u.DB, spotifyID.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}


// @Summary Deletes a user account by spotify ID
// @Description Deletes a user account by spotify ID
// @Tags Users
// @Accept json
// @Produce json
// @Success 204
// @Param spotifyID path string true "Spotify ID of the user to delete"
// @Failure 401 {string} string 
// @Failure 403 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/admin/{spotifyID} [delete]
// @Security Bearer
func(u *UserService) DeleteUserBySpotifyID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")

	err := u.UsersDAO.DeleteUserByID(u.DB, spotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func(u *UserService) updateUser(spotifyID string, userUpdateRequest *requests.UpdateUserRequestDTO) (*responses.User, error) {
	return u.UsersDAO.UpdateUserPropertiesBySpotifyID(u.DB, spotifyID, userUpdateRequest)
}
