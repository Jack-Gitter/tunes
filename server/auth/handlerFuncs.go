package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/server/auth/spotifyHelpers"
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

func GenerateJWT(c *gin.Context) {

    accessTokenResponse := spotifyHelpers.RetrieveAccessToken(c.Query("code"))
    userProfileResponse := spotifyHelpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    _, err := db.GetUserFromDbBySpotifyID(userProfileResponse.Id)

    if err != nil {
        err = db.InsertUserIntoDB(userProfileResponse.Id, userProfileResponse.Display_name, "user")
    }

    if err != nil {
        panic(err)
    }

    claims := &JWTClaims{
        RegisteredClaims: jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
        },
        SpotifyID: userProfileResponse.Id,
        AccessToken: accessTokenResponse.Access_token,
        RefreshToken: accessTokenResponse.Refresh_token,
        AccessTokenExpiresAt: accessTokenResponse.Expires_in,
        UserRole: "user",
    }

    claimsForRefresh := &jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsForRefresh)
    refreshString, _ := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.SetCookie("REFRESH_JWT", refreshString, 3600, "/", "localhost", false, true)

    c.Status(http.StatusOK)
}

type JWTClaims struct {
    SpotifyID string
    AccessToken string
    RefreshToken string
    AccessTokenExpiresAt int
    UserRole string
    jwt.RegisteredClaims
}

func (c JWTClaims) Validate() error {
    if c.SpotifyID != "blah" {
        // here, check if the access token is expeired for the user
        return errors.New("if you are seeing this, we still need to implement custom JWT claim validation")
    }
    return nil
}
