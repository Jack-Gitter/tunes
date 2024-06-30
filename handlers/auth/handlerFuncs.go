package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
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

    accessTokenResponse := retrieveAccessToken(c.Query("code"))
    userProfileResponse := retrieveUserProfile(accessTokenResponse.Access_token)

    claims :=  jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(), 
		"iat": time.Now().Unix(),
        "spotifyID": userProfileResponse.Id,
        "accessToken": accessTokenResponse.Access_token,
        "refreshToken": accessTokenResponse.Refresh_token,
        "accessTokenExpiresAt": accessTokenResponse.Expires_in,
        "userRole": "user",
    }        

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte("yadda"))

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.Status(http.StatusOK)
}

func retrieveAccessToken(authorizationCode string) *AccessTokenResponnse {

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

    accessTokenResponseBody := &AccessTokenResponnse{}
    json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

    return accessTokenResponseBody
}

func retrieveUserProfile(accessToken string) *ProfileResponse {

    nReq, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", &bytes.Buffer{})
    nReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    
    client := &http.Client{}
    nResp, _ := client.Do(nReq)

    respJson2 := &ProfileResponse{}

    json.NewDecoder(nResp.Body).Decode(respJson2)

    return respJson2

}

