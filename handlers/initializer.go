package handlers

import (
	"github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/login", auth.Login)
    r.GET("/generateJWT", auth.GenerateJWT)
    return r
}


