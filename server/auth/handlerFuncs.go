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
    }

    userProfileResponse, err := helpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to fetch the profile for the user")
    }

    _, err = db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        err = db.InsertUserIntoDB(userProfileResponse.Id, userProfileResponse.Display_name, "user")
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to create the new user")
    }

    tokenString, err := helpers.CreateAccessJWT(userProfileResponse.Id, accessTokenResponse.Access_token, accessTokenResponse.Refresh_token, accessTokenResponse.Expires_in)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to create a JWT access token for the user")
    }

    refreshString, err := helpers.CreateRefreshJWT()

    if err != nil {
        c.JSON(http.StatusInternalServerError, "unable to create a JWT refresh token for the user")
    }

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.SetCookie("REFRESH_JWT", refreshString, 3600, "/", "localhost", false, true)

    c.JSON(http.StatusOK, "login success")
}

func ValidateUserJWT(c *gin.Context) {
    
    jwtCookie, err := c.Cookie("JWT")

    if err != nil {
        c.JSON(http.StatusBadRequest, "no JWT access token provided. Please sign in before accessing this endpoint") 
    }

    token, err := helpers.ValidateAccessToken(jwtCookie)

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            spotifyID := token.Claims.(*models.JWTClaims).SpotifyID
            spotifyRefreshToken := token.Claims.(*models.JWTClaims).RefreshToken
            refreshJWT(c, spotifyID, spotifyRefreshToken)
        } else {
            c.JSON(http.StatusBadRequest, "nice try kid, don't fuck with the JWT")
        }
    }
}

func refreshJWT(c *gin.Context, spotifyID string, spotifyRefreshToken string) {

    refreshToken, err := c.Cookie("REFRESH_JWT")

    if err != nil {
        panic(err)
    }

    _, e := helpers.ValidateRefreshToken(refreshToken)

    if e != nil {
        if errors.Is(e, jwt.ErrTokenExpired) {
            c.JSON(http.StatusUnauthorized, "your refresh token has expired, please log back in")
        } else {
            c.JSON(http.StatusUnauthorized, "do not tamper with the refresh token either :) ")
        }
    }

    accessTokenResponseBody, err := helpers.RetreiveAccessTokenFromRefreshToken(refreshToken)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "error retreving a new spotify access token for the user")
    }

    if accessTokenResponseBody.Refresh_token == "" {
        accessTokenResponseBody.Refresh_token = spotifyRefreshToken
    }

    accessTokenJWT, err := helpers.CreateAccessJWT(spotifyID, accessTokenResponseBody.Access_token, accessTokenResponseBody.Refresh_token, accessTokenResponseBody.Expires_in)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "error creating a JWT for the user")
    }

    c.SetCookie("JWT", accessTokenJWT, 3600, "/", "localhost", false, true)
    //c.Next() ideally we just run this here and continue on with the user request, and then since we set the cookie they get it eventually
    c.JSON(http.StatusUnauthorized, "please make the request again, I have refreshed the token!!!")
}
