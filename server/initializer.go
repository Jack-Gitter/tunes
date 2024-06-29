package server

import (
	//"bytes"
	//"encoding/json"
	//	"fmt"
	"bytes"
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
    r.GET("/test", helloWorld)
    r.GET("/auth", auth)
    r.GET("/accesstoken", accessToken)
    return r
}

func helloWorld(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "hello world!",
    })
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
    client.Do(req) 

    claims :=  jwt.MapClaims{
		"sub": "test",
		"iss": "todo-app",                 
		"aud": "role",
		"exp": time.Now().Add(time.Hour).Unix(), 
		"iat": time.Now().Unix(),
        "spotifyID": "testID",
        "accessToken": "testToken",
        "refreshToken": "testRefresh",
        "accessTokenExpiresAt": "testDate",
    }        


    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte("yadda"))

    if err != nil {
        fmt.Println(err)
    }

    c.SetCookie("test_cookie", tokenString, 3600, "/", "localhost", false, true)
    
    c.JSON(200, gin.H{
        "message": "nice!",
    })
}



type AppClaims struct {
    spotifyID string 
    accessToken string 
    refreshToken string
    accessTokenExpiresAt string
    jwt.RegisteredClaims
}
