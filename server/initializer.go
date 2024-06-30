package server

import (
	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/login", auth.Login)
    r.GET("/generateJWT", auth.GenerateJWT)
    return r
}


