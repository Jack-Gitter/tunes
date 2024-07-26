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

    baseGroup := r.Group("", customerrors.ErrorHandlerMiddleware)
    {
        loginGroup := baseGroup.Group("/login") 
        {
            loginGroup.GET("/", auth.Login)
            loginGroup.GET("/callback", auth.LoginCallback)
            loginGroup.GET("/jwt", auth.RefreshJWT)
        }

        authGroup := baseGroup.Group("", auth.ValidateUserJWT) 
        {

            userGroup := authGroup.Group("/users")
            {
                userGroup.GET("/:spotifyID", users.GetUserById)
                userGroup.GET("/current", users.GetCurrentUser)
                userGroup.GET("/current/followers", users.GetFollowers)
                userGroup.GET("/:spotifyID/followers", users.GetFollowersByID)
                userGroup.POST("/current/follow/:otherUserSpotifyID", users.FollowerUser)
                userGroup.POST("/current/unfollow/:otherUserSpotifyID", users.UnFollowUser)
                userGroup.PATCH("/current", users.UpdateCurrentUserProperties)
                userGroup.DELETE("/current", users.DeleteCurrentUser)

                adminOnly := userGroup.Group("", auth.ValidateAdminUser)
                {
                    adminOnly.PATCH("/:spotifyID",  users.UpdateUserBySpotifyID)
                    adminOnly.DELETE("/:spotifyID", users.DeleteUserBySpotifyID)
                }

            }

            postGroup := authGroup.Group("/posts")
            {

                postGroup.GET("/:spotifyID/:songID", posts.GetPostBySpotifyIDAndSongID)
                postGroup.GET("/current/:songID", posts.GetPostCurrentUserBySongID)
                postGroup.GET("/previews/users/current", posts.GetAllPostsForCurrentUser)
                postGroup.GET("/previews/users/:spotifyID", posts.GetAllPostsForUserByID)
                postGroup.POST("/", posts.CreatePostForCurrentUser)
                postGroup.POST("/likes/:spotifyID/:songID", posts.LikePost)
                postGroup.POST("/dislikes/:spotifyID/:songID", posts.DislikePost)
                postGroup.PATCH("/current/:songID", posts.UpdateCurrentUserPost)
                postGroup.DELETE("/current/:songID", posts.DeletePostForCurrentUserBySongID)
                postGroup.DELETE("/votes/current/:posterSpotifyID/:songID",  posts.RemovePostVote)

                adminOnly := postGroup.Group("", auth.ValidateAdminUser)
                {
                    adminOnly.DELETE("/:spotifyID/:songID", posts.DeletePostBySpotifyIDAndSongID)
                }

            }
        }
    }

	return r
}
