package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	//"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/server/auth/spotifyHelpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
//	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Login(c *gin.Context) {

    client_id := os.Getenv("CLIENT_ID")
    scope := os.Getenv("SCOPES")
    redirect_uri := os.Getenv("REDIRECT_URI")

    endpoint := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s", client_id, scope, redirect_uri)

    c.Redirect(http.StatusMovedPermanently, endpoint) 
}

func GenerateJWT(c *gin.Context) {

    accessTokenResponse := spotifyHelpers.RetrieveAccessToken(c.Query("code"))
    userProfileResponse := spotifyHelpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    user, err := db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        user, err = db.InsertUserIntoDB(userProfileResponse.Id, userProfileResponse.Display_name, "user")
    }

    if err != nil {
        panic(err)
    }

    claims :=  jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(), 
		"iat": time.Now().Unix(),
        "spotifyID": userProfileResponse.Id,
        "accessToken": accessTokenResponse.Access_token,
        "refreshToken": accessTokenResponse.Refresh_token,
        "accessTokenExpiresAt": accessTokenResponse.Expires_in,
        "userRole": user.Role,
    }        

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte("yadda"))

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.Status(http.StatusOK)
}

