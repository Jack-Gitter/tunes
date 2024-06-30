package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/login", login)
    r.GET("/accesstoken", generateJWT)
    return r
}


func login(c *gin.Context) {

    client_id := "83ada5f0555a4f57be4243c3788cc9f4"
    scope := "user-read-private%20user-read-email"
    redirect_uri := "http://localhost:2000/accesstoken"

    endpoint := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s", client_id, scope, redirect_uri)

    c.Redirect(http.StatusMovedPermanently, endpoint) 
}

func retrieveAccessToken(authorizationCode string) *AccessTokenResponnse {

    queryParamsMap := url.Values{}
    queryParamsMap.Add("grant_type", "authorization_code")
    queryParamsMap.Add("code", authorizationCode)
    queryParamsMap.Add("redirect_uri", "http://localhost:2000/accesstoken")
    queryParams := queryParamsMap.Encode()

    req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(queryParams)))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", "Basic ODNhZGE1ZjA1NTVhNGY1N2JlNDI0M2MzNzg4Y2M5ZjQ6ODVkYjM4OTFiYzU0NDNmZGIxZDM2MDRjZjM5YTBhM2I=" )

    client := &http.Client{}
    resp, _ := client.Do(req) 
    respJson := &AccessTokenResponnse{}
    json.NewDecoder(resp.Body).Decode(respJson)
    return respJson
}

func retrieveUserProfile(accessToken string) *ProfileResponse {

    nReq, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", bytes.NewBuffer([]byte{}))
    nReq.Header.Set("Authorization", "Bearer " + accessToken) 
    
    client := &http.Client{}
    nResp, _ := client.Do(nReq)

    respJson2 := &ProfileResponse{}

    json.NewDecoder(nResp.Body).Decode(respJson2)

    return respJson2

}

func generateJWT(c *gin.Context) {

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



type AccessTokenResponnse struct {
    Access_token string 
    Token_type string 
    Scope string 
    Expires_in int 
    Refresh_token string 
}

type ProfileResponse struct {
    Id string
}

