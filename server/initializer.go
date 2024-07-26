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

    // all user endpoints tested and good!
    r.GET("/users/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetUserById)
    r.GET("/currentUser", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetCurrentUser)
    r.GET("/currentUser/followers", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetFollowers)
    r.GET("/users/:spotifyID/followers/", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.GetFollowersByID)
    r.POST("/currentUser/follow/:otherUserSpotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT,  users.FollowerUser)
    r.POST("/currentUser/unfollow/:otherUserSpotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.UnFollowUser)
    r.PATCH("/currentUser", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.UpdateCurrentUserProperties) 
    r.PATCH("/user/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, auth.ValidateAdminUser, users.UpdateUserBySpotifyID)
    r.DELETE("/currentUser", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, users.DeleteCurrentUser) 
    r.DELETE("/user/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, auth.ValidateAdminUser, users.DeleteUserBySpotifyID) 


    r.GET("/posts/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetPostBySpotifyIDAndSongID)
    r.GET("/currentUserPosts/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetPostCurrentUserBySongID)

    // fix pagniation
    r.GET("/currentUserPostPreviews/", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetAllPostsForCurrentUser)
    r.GET("/specificUserPostPreviews/:spotifyID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.GetAllPostsForUserByID)
    r.POST("/posts", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.CreatePostForCurrentUser)
    r.POST("/posts/like/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.LikePost)
    r.POST("/posts/dislike/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.DislikePost)
    r.PATCH("/currentUserPosts/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, posts.UpdateCurrentUserPost)
    r.DELETE("/posts/:spotifyID/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT, auth.ValidateAdminUser, posts.DeletePostBySpotifyIDAndSongID)
    r.DELETE("/currentUserPosts/:songID", customerrors.ErrorHandlerMiddleware, auth.ValidateUserJWT,  posts.DeletePostForCurrentUserBySongID)

    return r
}


