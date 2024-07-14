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

    r.GET("/users/:spotifyID", auth.ValidateUserJWT, auth.RefreshJWT, users.GetUserById)
    r.GET("/currentUser", auth.ValidateUserJWT, auth.RefreshJWT, users.GetCurrentUser)

    r.GET("/posts/:spotifyID/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.GetPostBySpotifyIDAndSongID)
    r.GET("/currentUserPosts/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.GetPostCurrentUserBySongID)
    r.POST("/posts", auth.ValidateUserJWT, auth.RefreshJWT, posts.CreatePostForCurrentUser)
    r.DELETE("/posts/:spotifyID/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.DeletePostBySpotifyIDAndSongID)
    r.DELETE("/currentUserPosts/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.DeletePostForCurrentUserBySongID)
    return r
}


