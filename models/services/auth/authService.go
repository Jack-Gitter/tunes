package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/daos"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/Jack-Gitter/tunes/models/services/jwt"
	"github.com/Jack-Gitter/tunes/models/services/spotify"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
    DB *sql.DB
    UsersDAO daos.IUsersDAO
    SpotifyService spotify.ISpotifyService
    JWTService jwt.IJWTService
}

type IAuthService interface {
    RefreshJWT(c *gin.Context) 
    ValidateUserJWT(c *gin.Context) 
    ValidateAdminUser(c *gin.Context) 
    Login(c *gin.Context) 
    LoginCallback(c *gin.Context) 
}

func(a *AuthService) Login(c *gin.Context) {

	client_id := os.Getenv("CLIENT_ID")
	scope := os.Getenv("SCOPES")
	redirect_uri := os.Getenv("REDIRECT_URI")

	endpoint := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s", client_id, scope, redirect_uri)

	c.Redirect(http.StatusMovedPermanently, endpoint)
}

func(a *AuthService) LoginCallback(c *gin.Context) {

	accessTokenResponse, err := a.SpotifyService.RetrieveInitialAccessToken(c.Query("code"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	userProfileResponse, err := a.SpotifyService.RetrieveUserProfile(accessTokenResponse.Access_token)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	user, err := a.UsersDAO.UpsertUserOnLogin(a.DB, userProfileResponse.Display_name, userProfileResponse.Id)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	tokenString, err := a.JWTService.CreateAccessJWT(
		userProfileResponse.Id,
		userProfileResponse.Display_name,
		accessTokenResponse.Access_token,
		accessTokenResponse.Expires_in,
		user.Role)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	refreshString, err := a.JWTService.CreateRefreshJWT(accessTokenResponse.Refresh_token)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.SetCookie("ACCESS_JWT", tokenString, 3600, "/", "localhost", false, false)
	c.SetCookie("REFRESH_JWT", refreshString, 86400, "/", "localhost", false, true)

	c.JSON(http.StatusOK, user)
}

// @Summary Refreshes the current users JWT
// @Description Refreshes the current users JWT
// @Tags Users
// @Accept json
// @Produce json
// @Param Cookie header string false "refresh JWT provided by login endpoint REFRESH_JWT=..."
// @Param spotifyID path string true "User Spotify ID"
// @Success 200 {object} responses.User 
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /users/{spotifyID} [get]
// @Security Bearer
func(a *AuthService) RefreshJWT(c *gin.Context) {

	refresh_jwt, err := c.Cookie("REFRESH_JWT")

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	refresh_token, e := a.JWTService.ValidateRefreshToken(refresh_jwt)

	if e != nil {
		c.Error(e)
		c.Abort()
		return
	}

	spotifyRefreshToken := refresh_token.Claims.(*requests.RefreshJWTClaims).RefreshToken
	accessTokenResponseBody, err := a.SpotifyService.RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if accessTokenResponseBody.Refresh_token == "" {
		accessTokenResponseBody.Refresh_token = spotifyRefreshToken
	}

	userProfileResponse, err := a.SpotifyService.RetrieveUserProfile(accessTokenResponseBody.Access_token)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	userDBResponse, err := a.UsersDAO.GetUserFromDbBySpotifyID(a.DB, userProfileResponse.Id)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	accessTokenJWT, err := a.JWTService.CreateAccessJWT(
		userProfileResponse.Id,
		userProfileResponse.Display_name,
		accessTokenResponseBody.Access_token,
		accessTokenResponseBody.Expires_in,
		userDBResponse.Role,
	)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.SetCookie("ACCESS_JWT", accessTokenJWT, 3600, "/", "localhost", false, false)

	c.Status(http.StatusNoContent)
}

func(a *AuthService) ValidateUserJWT(c *gin.Context) {

	header := strings.Split(c.GetHeader("Authorization"), " ")
	if len(header) < 2 {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusUnauthorized, Msg: "not enough values in the auth header"})
		c.Abort()
		return
	}

	if strings.ToLower(header[0]) != "bearer" {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid auth type"})
		c.Abort()
		return
	}

	jwtTokenString := header[1]

	token, err := a.JWTService.ValidateAccessToken(jwtTokenString)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	spotifyID := token.Claims.(*requests.JWTClaims).SpotifyID
	spotifyAccessToken := token.Claims.(*requests.JWTClaims).AccessToken
	role := token.Claims.(*requests.JWTClaims).UserRole
	username := token.Claims.(*requests.JWTClaims).Username

	c.Set("spotifyID", spotifyID)
	c.Set("userRole", role)
	c.Set("spotifyUsername", username)
	c.Set("spotifyAccessToken", spotifyAccessToken)

	c.Next()
}


func(a *AuthService) ValidateAdminUser(c *gin.Context) {

	role, found := c.Get("userRole")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "no role specified for user in db"})
		c.Abort()
		return
	}

	if role == responses.ADMIN {
		c.Next()
		return
	}

	c.Error(&customerrors.CustomError{StatusCode: http.StatusForbidden, Msg: "This endpoint is accessable to only admins!"})
	c.Abort()
	return

}
