package server

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"net/http"
	//"net/url"

	"github.com/gin-gonic/gin"
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
    //scope := "user-read-private user-read-email"
    //redirectUri := "http://localhost:80"
    endpoint := "https://accounts.spotify.com/authorize?response_type=code&client_id=83ada5f0555a4f57be4243c3788cc9f4&scope=user-read-private&redirect_uri=http://localhost:2000/accesstoken"

    c.Redirect(http.StatusMovedPermanently, endpoint) 

}

func accessToken(c *gin.Context) {
    fmt.Println(c.Query("code"))

   /* body := url.Values{}
    body.Set("grant_type", "authorization_code")
    body.Set("code", c.Query("code"))
    body.Set("redirect_uri", "http://localhost:2000/accesstoken")

    jsonString, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBuffer(jsonString))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", "Basic ODNhZGE1ZjA1NTVhNGY1N2JlNDI0M2MzNzg4Y2M5ZjQ6ODVkYjM4OTFiYzU0NDNmZGIxZDM2MDRjZjM5YTBhM2I=" )

    client := &http.Client{}
    resp, _ := client.Do(req) 
    fmt.Println(resp)*/
    c.Redirect(http.StatusFound, "http://localhost:8080")




    // for now, just return the token to the user. But we are going to want to wrap it in a jwt at some point with extra user information
    //req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token")
    //http.Post("https://accounts.spotify.com/api/token")
   
}



