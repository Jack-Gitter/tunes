package server

import (
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/validation"
	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/Jack-Gitter/tunes/server/comments"
	"github.com/Jack-Gitter/tunes/server/posts"
	"github.com/Jack-Gitter/tunes/server/users"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/Jack-Gitter/tunes/docs"
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
                userGroup.DELETE("/current/unfollow/:otherUserSpotifyID", users.UnFollowUser)

                // should not be able to update your role if you are not admin!
                userGroup.PATCH("/current", validation.ValidateData(requests.ValidateUserRequestDTO), users.UpdateCurrentUserProperties)
                userGroup.DELETE("/current", users.DeleteCurrentUser)

                adminOnly := userGroup.Group("/admin", auth.ValidateAdminUser)
                {
                    adminOnly.PATCH("/:spotifyID", validation.ValidateData(requests.ValidateUserRequestDTO), users.UpdateUserBySpotifyID)
                    adminOnly.DELETE("/:spotifyID", users.DeleteUserBySpotifyID)
                }

            }

            postGroup := authGroup.Group("/posts")
            {

                postGroup.GET("/:spotifyID/:songID", posts.GetPostBySpotifyIDAndSongID)
                postGroup.GET("/current/:songID", posts.GetPostCurrentUserBySongID)
                postGroup.GET("/previews/users/current", posts.GetAllPostsForCurrentUser)
                postGroup.GET("/previews/users/:spotifyID", posts.GetAllPostsForUserByID)
                postGroup.GET("/comments/:spotifyID/:songID", posts.GetPostCommentsPaginated)
                postGroup.POST("/", validation.ValidateData(requests.ValidateCreatePostDTO),  posts.CreatePostForCurrentUser)
                // should not be able to do so multiple times
                postGroup.POST("/likes/:spotifyID/:songID", posts.LikePost)
                postGroup.POST("/dislikes/:spotifyID/:songID", posts.DislikePost)
                // the likes do not get returned here... :(
                postGroup.PATCH("/current/:songID", validation.ValidateData(requests.ValidateUpdatePostRequestDTO), posts.UpdateCurrentUserPost)
                postGroup.DELETE("/current/:songID", posts.DeletePostForCurrentUserBySongID)

                postGroup.DELETE("/votes/current/:posterSpotifyID/:songID",  posts.RemovePostVote)

                adminOnly := postGroup.Group("/admin", auth.ValidateAdminUser)
                {
                    adminOnly.DELETE("/:spotifyID/:songID", posts.DeletePostBySpotifyIDAndSongID)
                }

            }

            commentGroup := authGroup.Group("/comments")
            {

                commentGroup.GET("/:commentID", comments.GetComment)
                commentGroup.POST("/:spotifyID/:songID", validation.ValidateData[requests.CreateCommentDTO](), comments.CreateComment)
                commentGroup.POST("/like/:commentID", comments.LikeComment)
                commentGroup.POST("/dislike/:commentID", comments.DislikeComment)
                commentGroup.PATCH("/current/:commentID", validation.ValidateData(requests.ValidateUpdateCommentDTO), comments.UpdateComment)
                commentGroup.DELETE("/current/:commentID", comments.DeleteCurrentUserComment)
                commentGroup.DELETE("/votes/current/:commentID", comments.RemoveCommentVote)

                adminOnly := commentGroup.Group("/admin", auth.ValidateAdminUser)
                {
                    adminOnly.DELETE("/:commentID", comments.DeleteComment)
                }

            }
        }
    }

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
