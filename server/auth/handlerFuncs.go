package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
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
        c.JSON(http.StatusInternalServerError, "unable to fetch the access token for the user from spotify")
        return
    }

    userProfileResponse, err := helpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to fetch the profile for the user")
        return
    }

    _, err = db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        err = db.InsertUserIntoDB(userProfileResponse.Id, userProfileResponse.Display_name, "user")
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to create the new user")
        return
    }

    tokenString, err := helpers.CreateAccessJWT(userProfileResponse.Id, userProfileResponse.Display_name, accessTokenResponse.Access_token, accessTokenResponse.Refresh_token, accessTokenResponse.Expires_in)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to create a JWT access token for the user")
        return
    }

    refreshString, err := helpers.CreateRefreshJWT()

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to create a JWT refresh token for the user")
        return
    }

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.SetCookie("REFRESH_JWT", refreshString, 3600, "/", "localhost", false, true)

    c.JSON(http.StatusOK, "login success")
}

func ValidateUserJWT(c *gin.Context) {
    
    jwtCookie, err := c.Cookie("JWT")

    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, "no JWT access token provided. Please sign in before accessing this endpoint") 
        return
    }

    token, err := helpers.ValidateAccessToken(jwtCookie)

    spotifyID := token.Claims.(*models.JWTClaims).SpotifyID
    spotifyRefreshToken := token.Claims.(*models.JWTClaims).RefreshToken
    spotifyUsername := token.Claims.(*models.JWTClaims).Username
    spotifyAccessToken := token.Claims.(*models.JWTClaims).AccessToken

    c.Set("spotifyID", spotifyID)
    c.Set("spotifyUsername", spotifyUsername)
    c.Set("spotifyAccessToken", spotifyAccessToken)
    c.Set("spotifyRefreshToken", spotifyRefreshToken)
    
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            c.Set("JWT_EXPIRED", true)
        } else {
            c.AbortWithStatusJSON(http.StatusBadRequest, "nice try kid, don't fuck with the JWT")
            return
        }
    } 
    c.Next()
}

func RefreshJWT(c *gin.Context) {

    _, exists := c.Get("JWT_EXPIRED")

    if !exists {
        c.Next()
    }

    spotifyID, _ := c.Get("spotifyID")
    spotifyUsername, _ := c.Get("spotifyUsername")
    spotifyRefreshToken, _ := c.Get("spotifyRefreshToken")


    refreshToken, err := c.Cookie("REFRESH_JWT")

    if err != nil {
        panic(err)
    }

    _, e := helpers.ValidateRefreshToken(refreshToken)

    if e != nil {
        if errors.Is(e, jwt.ErrTokenExpired) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "your refresh token has expired, please log back in")
            return
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "do not tamper with the refresh token either :) ")
            return
        }
    }

    accessTokenResponseBody, err := helpers.RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken.(string))

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "error retreiving a new spotify access token for the user")
        return
    }

    if accessTokenResponseBody.Refresh_token == "" {
        accessTokenResponseBody.Refresh_token = spotifyRefreshToken.(string)
    }

    accessTokenJWT, err := helpers.CreateAccessJWT(spotifyID.(string), spotifyUsername.(string), accessTokenResponseBody.Access_token, accessTokenResponseBody.Refresh_token, accessTokenResponseBody.Expires_in)

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "error creating a JWT for the user")
        return
    }

    c.SetCookie("JWT", accessTokenJWT, 3600, "/", "localhost", false, true)
    c.Set("spotifyUsername", spotifyUsername)
    c.Set("spotifyAccessToken", accessTokenResponseBody.Access_token)
    c.Next()
}
