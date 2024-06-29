package server

import (
	//"bytes"
	//"encoding/json"
	//	"fmt"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	//"io"
	//"time"

	//"encoding/json"
	//"fmt"
	"net/http"

	//"net/url"

	//"net/url"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/auth", auth)
    r.GET("/accesstoken", accessToken)
    return r
}


func auth(c *gin.Context) {
    endpoint := "https://accounts.spotify.com/authorize?response_type=code&client_id=83ada5f0555a4f57be4243c3788cc9f4&scope=user-read-private&redirect_uri=http://localhost:2000/accesstoken"
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

    fmt.Println("the json is!!!")
    fmt.Println(respJson)

    claims :=  jwt.MapClaims{
		"sub": "test",
		"iss": "todo-app",                 
		"aud": "role",
		"exp": time.Now().Add(time.Hour).Unix(), 
		"iat": time.Now().Unix(),
        "spotifyID": "testID",
        "accessToken": respJson.Access_token,
        "refreshToken": respJson.Refresh_token,
        "accessTokenExpiresAt": respJson.Expires_in,
    }        


    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, _ := token.SignedString([]byte("yadda"))

    c.SetCookie("test_cookie", tokenString, 3600, "/", "localhost", false, true)
    
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


