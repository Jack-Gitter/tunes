package server

import (
    "github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/test", helloWorld)
    return r
}

func helloWorld(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "hello world!",
    })
}



