package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Jack-Gitter/tunes/models"
	"github.com/golang-jwt/jwt/v5"
)

func RetrieveAccessToken(authorizationCode string) *models.AccessTokenResponnse {

    queryParamsMap := url.Values{}
    queryParamsMap.Add("grant_type", "authorization_code")
    queryParamsMap.Add("code", authorizationCode)
    queryParamsMap.Add("redirect_uri", os.Getenv("REDIRECT_URI"))
    queryParams := queryParamsMap.Encode()

    basicAuthToken := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
    encodedBasicAuthToken := base64.StdEncoding.EncodeToString([]byte(basicAuthToken))

    accessTokenRequest, _ := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(queryParams)))
    accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    accessTokenRequest.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedBasicAuthToken))

    client := &http.Client{}
    resp, _ := client.Do(accessTokenRequest) 

    accessTokenResponseBody := &models.AccessTokenResponnse{}
    json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

    return accessTokenResponseBody
}

func RetrieveUserProfile(accessToken string) *models.ProfileResponse {

    nReq, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", &bytes.Buffer{})
    nReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    
    client := &http.Client{}
    nResp, _ := client.Do(nReq)

    respJson2 := &models.ProfileResponse{}

    json.NewDecoder(nResp.Body).Decode(respJson2)

    return respJson2

}

func CreateAccessJWT(spotifyID string, accessToken string, refreshToken string, accessTokenExpiresAt int) (string, error) {

    claims := &models.JWTClaims{
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
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        AccessTokenExpiresAt: accessTokenExpiresAt,
        UserRole: "user",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

    return tokenString, err
}

func CreateRefreshJWT() (string, error) {

    claimsForRefresh := &jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsForRefresh)
    refreshString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

    return refreshString, err
}

func ValidateAccessToken(accessTokenJWT string) (*jwt.Token, error) {
    token, err := jwt.ParseWithClaims(accessTokenJWT, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    return token, err
}

func ValidateRefreshToken(refreshTokenJWT string) (*jwt.Token, error) {
    token, e := jwt.Parse(refreshTokenJWT, func (token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    return token, e
}
