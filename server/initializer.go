package server

import (
	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/Jack-Gitter/tunes/server/users"
	"github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
    r := gin.Default()
    r.GET("/login", auth.Login)
    r.GET("/loginCallback", auth.LoginCallback)
    //r.GET("/validate", auth.ValidateUserJWT) // needs to be a middleware on endpoints
    r.GET("/user", auth.ValidateUserJWT, users.GetUserById)
    return r
}


