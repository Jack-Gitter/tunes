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

    user, err := db.UpsertUser(userProfileResponse.Display_name, userProfileResponse.Id)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    tokenString, err := helpers.CreateAccessJWT(
        userProfileResponse.Id, 
        userProfileResponse.Display_name, 
        accessTokenResponse.Access_token, 
        accessTokenResponse.Expires_in, 
        user.Role) 

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    refreshString, err := helpers.CreateRefreshJWT(accessTokenResponse.Refresh_token)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    c.SetCookie("ACCESS_JWT", tokenString, 3600, "/", "localhost", false, false)
    c.SetCookie("REFRESH_JWT", refreshString, 86400, "/", "localhost", false, true)

    c.JSON(http.StatusOK, user)
}

func ValidateUserJWT(c *gin.Context) {
    
    header := strings.Split(c.GetHeader("Authorization"), " ")
    if len(header) < 2 {
        c.AbortWithStatusJSON(http.StatusBadRequest, "Not enough values supplied in the Authorization header")
        return
    }

    if strings.ToLower(header[0]) != "bearer" {
        c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid authorization type")
        return
    }

    jwtTokenString := header[1]

    token, err := helpers.ValidateAccessToken(jwtTokenString)

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "Please refresh your JWT by accessing the refresh endpoint")
            return
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "JWT signature could not be verified")
            return
        }
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
        c.JSON(http.StatusBadRequest, "Refresh JWT cookie missing")
        return
    }

    refresh_token, e := helpers.ValidateRefreshToken(refresh_jwt)

    if e != nil {
        if errors.Is(e, jwt.ErrTokenExpired) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "Your refresh JWT has expired, please log out and log back in to continue")
            return
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, "Refresh JWT signature could not be verified")
            return
        }
    }

    spotifyRefreshToken := refresh_token.Claims.(*requests.RefreshJWTClaims).RefreshToken
    accessTokenResponseBody, err := helpers.RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken)

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "Error retrieving a new spotify access token")
        return
    }

    if accessTokenResponseBody.Refresh_token == "" {
        accessTokenResponseBody.Refresh_token = spotifyRefreshToken
    }

    userProfileResponse, err := helpers.RetrieveUserProfile(accessTokenResponseBody.Access_token)

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
        return
    }

    userDBResponse, err := db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        if err, ok := err.(*db.DBError); ok {
            c.AbortWithStatusJSON(err.StatusCode, err.Msg)
            return
        }
        panic("cant be here")
    }

    accessTokenJWT, err := helpers.CreateAccessJWT(
        userProfileResponse.Id, 
        userProfileResponse.Display_name, 
        accessTokenResponseBody.Access_token, 
        accessTokenResponseBody.Expires_in,
        userDBResponse.Role,
    )

    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
        return
    }

    c.SetCookie("ACCESS_JWT", accessTokenJWT, 3600, "/", "localhost", false, false)

    c.Status(http.StatusNoContent)
}

func ValidateAdminUser(c *gin.Context) {

    role, found := c.Get("userRole")

    if !found {
        c.AbortWithStatusJSON(http.StatusInternalServerError, "No role specified for the current user in the database")
        return
    }

    if role == responses.ADMIN {
        c.Next()
        return
    }

    c.AbortWithStatusJSON(http.StatusForbidden, "Only admins are allowed to access this endpoint")
    return

}
