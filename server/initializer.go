package server

import (
	"bytes"
	"encoding/json"
	"fmt"

	//"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/login", login)
    r.GET("/accesstoken", accessToken)
    return r
}


func login(c *gin.Context) {
    endpoint := "https://accounts.spotify.com/authorize?response_type=code&client_id=83ada5f0555a4f57be4243c3788cc9f4&scope=user-read-private%20user-read-email&redirect_uri=http://localhost:2000/accesstoken"
    c.Redirect(http.StatusMovedPermanently, endpoint) 
}

func accessToken(c *gin.Context) {

    b := "grant_type=authorization_code&code=" + c.Query("code") + "&redirect_uri=http://localhost:2000/accesstoken"

    req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(b)))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", "Basic ODNhZGE1ZjA1NTVhNGY1N2JlNDI0M2MzNzg4Y2M5ZjQ6ODVkYjM4OTFiYzU0NDNmZGIxZDM2MDRjZjM5YTBhM2I=" )

    client := &http.Client{}
    resp, _ := client.Do(req) 

    respJson := &AccessTokenResponnse{}

    json.NewDecoder(resp.Body).Decode(respJson)


    nReq, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", bytes.NewBuffer([]byte{}))
    nReq.Header.Set("Authorization", "Bearer " + respJson.Access_token)
    
    nResp, _ := client.Do(nReq)

    //fmt.Println(nResp.Body)
    respJson2 := &ProfileResponse{}

    json.NewDecoder(nResp.Body).Decode(respJson2)

    fmt.Println(respJson2.Id)
    // fetch the users access token here

    claims :=  jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(), 
		"iat": time.Now().Unix(),
        "spotifyID": respJson2.Id,
        "accessToken": respJson.Access_token,
        "refreshToken": respJson.Refresh_token,
        "accessTokenExpiresAt": respJson.Expires_in,
        "userRole": "user",
    }        

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, _ := token.SignedString([]byte("yadda"))

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    
    c.JSON(200, gin.H{
        "message": "nice!",
    })
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

