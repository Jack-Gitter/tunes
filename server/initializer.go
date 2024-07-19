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
    r.GET("/currentUser/followers", auth.ValidateUserJWT, auth.RefreshJWT, users.GetFollowers)
    r.GET("/users/:spotifyID/followers/", auth.ValidateUserJWT, auth.RefreshJWT, users.GetFollowersByID)
    r.POST("/currentUser/follow/:otherUserSpotifyID", auth.ValidateUserJWT, auth.RefreshJWT, users.FollowerUser)
    r.POST("/currentUser/unfollow/:otherUserSpotifyID", auth.ValidateUserJWT, auth.RefreshJWT, users.UnFollowUser)
    r.PATCH("/currentUser", auth.ValidateUserJWT, auth.RefreshJWT, users.UpdateCurrentUserProperties) 
    r.PATCH("/user/:spotifyID", auth.ValidateUserJWT, auth.RefreshJWT, users.UpdateUserBySpotifyID)
    r.DELETE("/currentUser", auth.ValidateUserJWT, auth.RefreshJWT, users.DeleteCurrentUser) 
    r.DELETE("/user/:spotifyID", auth.ValidateUserJWT, auth.RefreshJWT, auth.ValidateAdminUser, users.DeleteUserBySpotifyID) 


    r.GET("/posts/:spotifyID/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.GetPostBySpotifyIDAndSongID)
    r.GET("/currentUserPosts/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.GetPostCurrentUserBySongID)
    r.GET("/currentUserPostPreviews/", auth.ValidateUserJWT, auth.RefreshJWT, posts.GetAllPostsForCurrentUser)
    r.GET("/specificUserPostPreviews/:spotifyID", auth.ValidateUserJWT, auth.RefreshJWT, posts.GetAllPostsForUserByID)
    r.POST("/posts", auth.ValidateUserJWT, auth.RefreshJWT, posts.CreatePostForCurrentUser)
    r.POST("/currentUserPosts/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.UpdateCurrentUserPost)
    r.POST("/posts/like/:spotifyID/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.LikePost)
    r.DELETE("/posts/:spotifyID/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.DeletePostBySpotifyIDAndSongID)
    r.DELETE("/currentUserPosts/:songID", auth.ValidateUserJWT, auth.RefreshJWT, posts.DeletePostForCurrentUserBySongID)

    return r
}


