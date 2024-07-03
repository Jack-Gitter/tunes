package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
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

func GenerateJWT(c *gin.Context) {

    accessTokenResponse := helpers.RetrieveInitialAccessToken(c.Query("code"))
    userProfileResponse := helpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    _, err := db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        err = db.InsertUserIntoDB(userProfileResponse.Id, userProfileResponse.Display_name, "user")
    }

    if err != nil {
        panic(err)
    }

    tokenString, err := helpers.CreateAccessJWT(userProfileResponse.Id, accessTokenResponse.Access_token, accessTokenResponse.Refresh_token, accessTokenResponse.Expires_in)

    if err != nil {
        panic(err)
    }

    refreshString, err := helpers.CreateRefreshJWT()

    if err != nil {
        panic(err)
    }

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.SetCookie("REFRESH_JWT", refreshString, 3600, "/", "localhost", false, true)

    c.Status(http.StatusOK)
}

func ValidateUserJWT(c *gin.Context) {
    
    jwtCookie, err := c.Cookie("JWT")

    if err != nil {
        panic(err)
    }

    token, err := helpers.ValidateAccessToken(jwtCookie)

    if err != nil {
        // here check for specific error?
        spotifyID := token.Claims.(*models.JWTClaims).SpotifyID
        spotifyRefreshToken := token.Claims.(*models.JWTClaims).RefreshToken
        refreshJWT(c, spotifyID, spotifyRefreshToken)
    }
}

func refreshJWT(c *gin.Context, spotifyID string, spotifyRefreshToken string) {

    refreshToken, err := c.Cookie("REFRESH_JWT")

    if err != nil {
        panic(err)
    }

    _, e := helpers.ValidateRefreshToken(refreshToken)

    if e != nil {
        c.JSON(http.StatusUnauthorized, "the refresh token has expired. Please log out and log back in again")
    }

    accessTokenResponseBody, err := helpers.RetreiveAccessTokenFromRefreshToken(refreshToken)

    if accessTokenResponseBody.Refresh_token == "" {
        accessTokenResponseBody.Refresh_token = spotifyRefreshToken
    }

    accessTokenJWT, err := helpers.CreateAccessJWT(spotifyID, accessTokenResponseBody.Access_token, accessTokenResponseBody.Refresh_token, accessTokenResponseBody.Expires_in)

    c.SetCookie("JWT", accessTokenJWT, 3600, "/", "localhost", false, true)
    //c.Next() ideally we just run this here and continue on with the user request, and then since we set the cookie they get it eventually
    c.JSON(http.StatusUnauthorized, "please make the request again, I have refreshed the token!!!")
}
