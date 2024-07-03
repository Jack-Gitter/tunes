package middlware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Jack-Gitter/tunes/server/auth"

	//"github.com/Jack-Gitter/tunes/server/auth/spotifyHelpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateUserJWT(c *gin.Context) {
    
    jwtCookie, err := c.Cookie("JWT")
    if err != nil {
        panic(err)
    }

    token, err := jwt.ParseWithClaims(jwtCookie, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        spotifyID := token.Claims.(*auth.JWTClaims).SpotifyID
        spotifyRefreshToken := token.Claims.(*auth.JWTClaims).RefreshToken
        refreshJWT(c, spotifyID, spotifyRefreshToken)
    }
}

func refreshJWT(c *gin.Context, spotifyID string, spotifyRefreshToken string) {

    refreshToken, err := c.Cookie("REFRESH_JWT")
    fmt.Println(spotifyID)

    if err != nil {
        panic(err)
    }

    _, e := jwt.Parse(refreshToken, func (token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if e != nil {
        // reroute the user to login screen, because refresh token has expired
        c.JSON(http.StatusUnauthorized, "the refresh token has expired. Please log out and log back in again")
    }

    // generate a new spotify access token, refresh token, and expires at and put them below
    queryParamsMap := url.Values{}
    queryParamsMap.Add("grant_type", "refresh_token")
    queryParamsMap.Add("refresh_token", spotifyRefreshToken)
    queryParams := queryParamsMap.Encode()

    basicAuthToken := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
    encodedBasicAuthToken := base64.StdEncoding.EncodeToString([]byte(basicAuthToken))

    accessTokenRefreshRequest, _ := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(queryParams)))
    accessTokenRefreshRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    accessTokenRefreshRequest.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedBasicAuthToken))

    client := &http.Client{}
    resp, _ := client.Do(accessTokenRefreshRequest) 


    accessTokenResponseBody := &RefreshTokenResponse{}
    json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

    if accessTokenResponseBody.Refresh_token == "" {
        accessTokenResponseBody.Refresh_token = spotifyRefreshToken
    }

    claims := &auth.JWTClaims{
        RegisteredClaims: jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           ExpiresAt: &jwt.NumericDate{Time: time.Now()},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
        },
        SpotifyID: spotifyID,
        AccessToken: accessTokenResponseBody.Access_token,
        RefreshToken: accessTokenResponseBody.Refresh_token,
        AccessTokenExpiresAt: accessTokenResponseBody.Expires_in,
        UserRole: "user",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.JSON(http.StatusUnauthorized, "please make the request again, I have refreshed the token!!!")
}

type RefreshTokenResponse struct {
    Access_token string
    Expires_in int
    Refresh_token string
}
