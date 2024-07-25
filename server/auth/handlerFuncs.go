package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/Jack-Gitter/tunes/server/auth/helpers"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

    client_id := os.Getenv("CLIENT_ID")
    scope := os.Getenv("SCOPES")
    redirect_uri := os.Getenv("REDIRECT_URI")

    endpoint := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s", client_id, scope, redirect_uri)

    c.Redirect(http.StatusMovedPermanently, endpoint) 
}

func LoginCallback(c *gin.Context) {

    accessTokenResponse, err := helpers.RetrieveInitialAccessToken(c.Query("code"))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    userProfileResponse, err := helpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    user, err := db.UpsertUser(userProfileResponse.Display_name, userProfileResponse.Id)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    tokenString, err := helpers.CreateAccessJWT(
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


    refreshString, err := helpers.CreateRefreshJWT(accessTokenResponse.Refresh_token)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.SetCookie("ACCESS_JWT", tokenString, 3600, "/", "localhost", false, false)
    c.SetCookie("REFRESH_JWT", refreshString, 86400, "/", "localhost", false, true)

    c.JSON(http.StatusOK, user)
}

func ValidateUserJWT(c *gin.Context) {
    
    header := strings.Split(c.GetHeader("Authorization"), " ")
    if len(header) < 2 {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "not enough values in the auth header"})
        c.Abort()
        return
    }

    if strings.ToLower(header[0]) != "bearer" {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid auth type"})
        c.Abort()
        return
    }

    jwtTokenString := header[1]

    token, err := helpers.ValidateAccessToken(jwtTokenString)

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

func RefreshJWT(c *gin.Context) {

    refresh_jwt, err := c.Cookie("REFRESH_JWT")

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    refresh_token, e := helpers.ValidateRefreshToken(refresh_jwt)

    if e != nil {
        c.Error(e)
        c.Abort()
        return
    }

    spotifyRefreshToken := refresh_token.Claims.(*requests.RefreshJWTClaims).RefreshToken
    accessTokenResponseBody, err := helpers.RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    if accessTokenResponseBody.Refresh_token == "" {
        accessTokenResponseBody.Refresh_token = spotifyRefreshToken
    }

    userProfileResponse, err := helpers.RetrieveUserProfile(accessTokenResponseBody.Access_token)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    userDBResponse, err := db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    accessTokenJWT, err := helpers.CreateAccessJWT(
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

func ValidateAdminUser(c *gin.Context) {

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

    c.Error(&customerrors.CustomError{StatusCode: http.StatusForbidden, Msg: "only admins"})
    c.Abort()
    return

}
