package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/Jack-Gitter/tunes/server/auth/helpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    userProfileResponse, err := helpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    user, err := db.UpsertUser(userProfileResponse.Display_name, userProfileResponse.Id, responses.BASIC_USER)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    tokenString, err := helpers.CreateAccessJWT(
        userProfileResponse.Id, 
        userProfileResponse.Display_name, 
        accessTokenResponse.Access_token, 
        accessTokenResponse.Refresh_token, 
        accessTokenResponse.Expires_in, 
        user.Role) 

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    refreshString, err := helpers.CreateRefreshJWT()

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    resp := responses.AuthResponse{}
    resp.AccessToken = tokenString
    resp.RefreshToken = refreshString
    resp.Data = user

    c.JSON(http.StatusOK, resp)
}

func ValidateUserJWT(c *gin.Context) {
    
    header := strings.Split(c.GetHeader("Authorization"), " ")
    if len(header) < 2 {
        c.AbortWithStatusJSON(http.StatusBadRequest, "invalid auth header!")
        return
    }

    if strings.ToLower(header[0]) != "bearer" {
        c.AbortWithStatusJSON(http.StatusBadRequest, "invalid auth type!")
        return
    }

    jwtTokenString := header[1]

    token, err := helpers.ValidateAccessToken(jwtTokenString)

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "please refresh your JWT with the provided refresh token!")
            return
        } else {
            c.AbortWithStatusJSON(http.StatusBadRequest, "nice try kid, don't fuck with the JWT")
            return
        }
    } 

    spotifyID := token.Claims.(*requests.JWTClaims).SpotifyID
    spotifyRefreshToken := token.Claims.(*requests.JWTClaims).RefreshToken
    spotifyAccessToken := token.Claims.(*requests.JWTClaims).AccessToken
    role := token.Claims.(*requests.JWTClaims).UserRole

    c.Set("spotifyID", spotifyID)
    c.Set("userRole", role)
    c.Set("spotifyAccessToken", spotifyAccessToken)
    c.Set("spotifyRefreshToken", spotifyRefreshToken)
    
    c.Next()
}

func RefreshJWT(c *gin.Context) {

    authReq := &requests.RefreshJWTDTO{}
    err := c.ShouldBindBodyWithJSON(authReq)

    if err != nil {
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "bad data for refresh")
        return
    }

    _, e := helpers.ValidateRefreshToken(authReq.RefreshToken)

    if e != nil {
        if errors.Is(e, jwt.ErrTokenExpired) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "your refresh token has expired, please log back in")
            return
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "do not tamper with the refresh token either :) ")
            return
        }
    }

    acc_token, e := helpers.ValidateAccessToken(authReq.AccessToken)

    if e != nil && !errors.Is(e, jwt.ErrTokenExpired) {
        c.AbortWithStatusJSON(http.StatusUnauthorized, "do not fuck with the JWT BITCH")
        return
    }

    spotifyRefreshToken := acc_token.Claims.(*requests.JWTClaims).RefreshToken
    accessTokenResponseBody, err := helpers.RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken)

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "error retreiving a new spotify access token for the user")
        return
    }

    if accessTokenResponseBody.Refresh_token == "" {
        accessTokenResponseBody.Refresh_token = spotifyRefreshToken
    }

    userProfileResponse, err := helpers.RetrieveUserProfile(accessTokenResponseBody.Access_token)

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "unable to get user profile from spotify")
        return
    }

    userDBResponse, found, err := db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "unable to get user from db")
        return
    }

    if !found {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "user does not exist in DB")
        return
    }

    accessTokenJWT, err := helpers.CreateAccessJWT(
        userProfileResponse.Id, 
        userProfileResponse.Display_name, 
        accessTokenResponseBody.Access_token, 
        accessTokenResponseBody.Refresh_token, 
        accessTokenResponseBody.Expires_in,
        userDBResponse.Role,
    )

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "error creating a JWT for the user")
        return
    }

    resp := map[string]string{
        "AccessToken": accessTokenJWT,
    }

    c.JSON(http.StatusOK, resp)
}

func ValidateAdminUser(c *gin.Context) {

    role, found := c.Get("userRole")

    if !found {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "could not get role for the current user!")
        return
    }

    if role == responses.ADMIN {
        c.Next()
        return
    }

    c.AbortWithStatusJSON(http.StatusBadRequest, "cannot access this endpoint if your not admin dummy!")
    return

}
