package auth

import (
	"fmt"
	"net/http"
	"os"
	"github.com/Jack-Gitter/tunes/db"
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

    accessTokenResponse := helpers.RetrieveAccessToken(c.Query("code"))
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
