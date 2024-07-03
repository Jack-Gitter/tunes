package server

import (
	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/Jack-Gitter/tunes/server/posts"
	"github.com/Jack-Gitter/tunes/server/users"
	"github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()

    r.GET("/login", auth.Login)
    r.GET("/loginCallback", auth.LoginCallback)

    r.GET("/user", auth.ValidateUserJWT, users.GetUserById)
    r.GET("/currentUser", auth.ValidateUserJWT, users.GetCurrentUser)

    r.POST("/post", auth.ValidateUserJWT, posts.CreatePostForCurrentUser)
    return r
}


