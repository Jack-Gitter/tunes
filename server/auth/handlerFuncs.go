package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"

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

    user, foundUser, err := db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !foundUser {
        user, err = db.InsertUserIntoDB(userProfileResponse.Display_name, userProfileResponse.Id, responses.BASIC_USER)
    }

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

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.SetCookie("REFRESH_JWT", refreshString, 3600, "/", "localhost", false, true)

    c.JSON(http.StatusOK, user)
}

func ValidateUserJWT(c *gin.Context) {
    
    jwtCookie, err := c.Cookie("JWT")

    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, "no JWT access token provided. Please sign in before accessing this endpoint") 
        return
    }

    token, err := helpers.ValidateAccessToken(jwtCookie)

    spotifyID := token.Claims.(*requests.JWTClaims).SpotifyID
    spotifyRefreshToken := token.Claims.(*requests.JWTClaims).RefreshToken
    spotifyAccessToken := token.Claims.(*requests.JWTClaims).AccessToken
    role := token.Claims.(*requests.JWTClaims).UserRole

    c.Set("spotifyID", spotifyID)
    c.Set("userRole", role)
    c.Set("spotifyAccessToken", spotifyAccessToken)
    c.Set("spotifyRefreshToken", spotifyRefreshToken)
    
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            fmt.Println("here")
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
        return
    }

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

    c.SetCookie("JWT", accessTokenJWT, 3600, "/", "localhost", false, true)
    c.Set("spotifyID", userProfileResponse.Id)
    c.Set("userRole", userDBResponse.Role)
    c.Set("spotifyAccessToken", accessTokenResponseBody.Access_token)
    c.Next()
}

func ValidateAdminUser(c *gin.Context) {

    role, found := c.Get("userRole")

    if !found {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "could not get role for the current user!")
        return
    }

    if role == responses.ADMIN {
        c.Next()
    }

    c.AbortWithStatusJSON(http.StatusBadRequest, "cannot access this endpoint if your not admin dummy!")
    return

}
