package server

import (
	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/Jack-Gitter/tunes/server/posts"
	"github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/login", auth.Login)
    r.GET("/loginCallback", auth.LoginCallback)
    r.POST("/post", posts.CreatePost)
    r.GET("/validate", auth.ValidateUserJWT)
    r.GET("/refreshJWT")
    return r
}


