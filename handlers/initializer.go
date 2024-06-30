package handlers

import (
	handlers "github.com/Jack-Gitter/tunes/handlers/auth"
	"github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/login", handlers.Login)
    r.GET("/generateJWT", handlers.GenerateJWT)
    return r
}


