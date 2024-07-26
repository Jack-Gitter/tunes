package server

import (
	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/Jack-Gitter/tunes/server/posts"
	"github.com/Jack-Gitter/tunes/server/users"
	"github.com/gin-gonic/gin"
)

func InitializeHttpServer() *gin.Engine {
	r := gin.Default()

	r.GET("/login", customerrors.ErrorHandlerMiddleware, auth.Login)
	r.GET("/loginCallback", customerrors.ErrorHandlerMiddleware, auth.LoginCallback)

	r.POST("/refreshJWT", customerrors.ErrorHandlerMiddleware, auth.RefreshJWT)

	r.GET("/users/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetUserById)
	r.GET("/users/current", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetCurrentUser)
	r.GET("/users/current/followers", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetFollowers)
	r.GET("/users/:spotifyID/followers", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetFollowersByID)
	r.POST("/users/current/follow/:otherUserSpotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.FollowerUser)
	r.POST("/users/current/unfollow/:otherUserSpotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.UnFollowUser)
	r.PATCH("/users/current", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.UpdateCurrentUserProperties)
	r.PATCH("/users/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, auth.ValidateAdminUser, users.UpdateUserBySpotifyID)
	r.DELETE("/users/current", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.DeleteCurrentUser)
    r.DELETE("/users/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, auth.ValidateAdminUser, users.DeleteUserBySpotifyID)

	r.GET("/posts/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetPostBySpotifyIDAndSongID)
	r.GET("/posts/current/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetPostCurrentUserBySongID)
	r.GET("/posts/previews/users/current", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetAllPostsForCurrentUser)
	r.GET("/posts/previews/users/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetAllPostsForUserByID)
	r.POST("/posts", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.CreatePostForCurrentUser)
	r.POST("/posts/likes/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.LikePost)
	r.POST("/posts/dislikes/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.DislikePost)
	r.PATCH("/posts/current/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.UpdateCurrentUserPost)
	r.DELETE("/posts/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, auth.ValidateAdminUser, posts.DeletePostBySpotifyIDAndSongID)
	r.DELETE("/posts/current/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.DeletePostForCurrentUserBySongID)
	r.DELETE("/posts/votes/current/:posterSpotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.RemovePostVote)

	return r
}
