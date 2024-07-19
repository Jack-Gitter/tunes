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

    r.POST("/refreshJWT", auth.RefreshJWT)

    r.GET("/users/:spotifyID", auth.ValidateUserJWT, users.GetUserById)
    r.GET("/currentUser", auth.ValidateUserJWT, users.GetCurrentUser)
    r.GET("/currentUser/followers", auth.ValidateUserJWT, users.GetFollowers)
    r.GET("/users/:spotifyID/followers/", auth.ValidateUserJWT, users.GetFollowersByID)
    r.POST("/currentUser/follow/:otherUserSpotifyID", auth.ValidateUserJWT,  users.FollowerUser)
    r.POST("/currentUser/unfollow/:otherUserSpotifyID", auth.ValidateUserJWT, users.UnFollowUser)
    r.PATCH("/currentUser", auth.ValidateUserJWT, users.UpdateCurrentUserProperties) 
    r.PATCH("/user/:spotifyID", auth.ValidateUserJWT, users.UpdateUserBySpotifyID)
    r.DELETE("/currentUser", auth.ValidateUserJWT, users.DeleteCurrentUser) 
    r.DELETE("/user/:spotifyID", auth.ValidateUserJWT, auth.ValidateAdminUser, users.DeleteUserBySpotifyID) 


    r.GET("/posts/:spotifyID/:songID", auth.ValidateUserJWT, posts.GetPostBySpotifyIDAndSongID)
    r.GET("/currentUserPosts/:songID", auth.ValidateUserJWT, posts.GetPostCurrentUserBySongID)
    r.GET("/currentUserPostPreviews/", auth.ValidateUserJWT, posts.GetAllPostsForCurrentUser)
    r.GET("/specificUserPostPreviews/:spotifyID", auth.ValidateUserJWT, posts.GetAllPostsForUserByID)
    r.POST("/posts", auth.ValidateUserJWT, posts.CreatePostForCurrentUser)
    r.POST("/currentUserPosts/:songID", auth.ValidateUserJWT, posts.UpdateCurrentUserPost)
    r.POST("/posts/like/:spotifyID/:songID", auth.ValidateUserJWT, posts.LikePost)
    r.DELETE("/posts/:spotifyID/:songID", auth.ValidateUserJWT,  posts.DeletePostBySpotifyIDAndSongID)
    r.DELETE("/currentUserPosts/:songID", auth.ValidateUserJWT,  posts.DeletePostForCurrentUserBySongID)

    return r
}


