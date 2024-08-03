package server

import (
	_ "github.com/Jack-Gitter/tunes/docs"
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
)

func InitializeHttpServer(userService users.IUserSerivce, postsService posts.IPostsService, commentsService comments.ICommentsService, authSerivce auth.IAuthService) *gin.Engine {
	r := gin.Default()

    baseGroup := r.Group("", customerrors.ErrorHandlerMiddleware)
    {
        loginGroup := baseGroup.Group("/login") 
        {
            loginGroup.GET("/", authSerivce.Login)
            loginGroup.GET("/callback", authSerivce.LoginCallback)
            loginGroup.GET("/jwt", authSerivce.RefreshJWT)
        }

        authGroup := baseGroup.Group("", authSerivce.ValidateUserJWT) 
        {

            userGroup := authGroup.Group("/users")
            {
                userGroup.GET("/:spotifyID", userService.GetUserById)
                userGroup.GET("/current", userService.GetCurrentUser)
                userGroup.GET("/current/followers", userService.GetFollowers)
                userGroup.GET("/:spotifyID/followers", userService.GetFollowersByID)
                userGroup.POST("/current/follow/:otherUserSpotifyID", userService.FollowUser)
                userGroup.DELETE("/current/unfollow/:otherUserSpotifyID", userService.UnFollowUser)
                userGroup.PATCH("/current", validation.ValidateData(requests.ValidateUserRequestDTO), userService.UpdateCurrentUserProperties)
                userGroup.DELETE("/current", userService.DeleteCurrentUser)

                adminOnly := userGroup.Group("/admin", authSerivce.ValidateAdminUser)
                {
                    adminOnly.PATCH("/:spotifyID", validation.ValidateData(requests.ValidateUserRequestDTO), userService.UpdateUserBySpotifyID)
                    adminOnly.DELETE("/:spotifyID", userService.DeleteUserBySpotifyID)
                }

            }

            postGroup := authGroup.Group("/posts")
            {

                postGroup.GET("/:spotifyID/:songID", postsService.GetPostBySpotifyIDAndSongID)
                postGroup.GET("/current/:songID", postsService.GetPostCurrentUserBySongID)
                postGroup.GET("/previews/users/current", postsService.GetAllPostsForCurrentUser)
                postGroup.GET("/previews/users/:spotifyID", postsService.GetAllPostsForUserByID)
                postGroup.GET("/comments/:spotifyID/:songID", postsService.GetPostCommentsPaginated)
                postGroup.GET("/feed", postsService.GetCurrentUserFeed)
                postGroup.POST("/", validation.ValidateData(requests.ValidateCreatePostDTO),  postsService.CreatePostForCurrentUser)
                postGroup.POST("/likes/:spotifyID/:songID", postsService.LikePost)
                postGroup.POST("/dislikes/:spotifyID/:songID", postsService.DislikePost)
                postGroup.PATCH("/current/:songID", validation.ValidateData(requests.ValidateUpdatePostRequestDTO), postsService.UpdateCurrentUserPost)
                postGroup.DELETE("/current/:songID", postsService.DeletePostForCurrentUserBySongID)
                postGroup.DELETE("/votes/current/:posterSpotifyID/:songID",  postsService.RemovePostVote)

                adminOnly := postGroup.Group("/admin", authSerivce.ValidateAdminUser)
                {
                    adminOnly.DELETE("/:spotifyID/:songID", postsService.DeletePostBySpotifyIDAndSongID)
                }

            }

            commentGroup := authGroup.Group("/comments")
            {

                commentGroup.GET("/:commentID",  validation.ValidatePathParams[requests.CommentIDPathParams](), commentsService.GetComment)
                commentGroup.POST("/:spotifyID/:songID", validation.ValidatePathParams[requests.CommentIDPathParams](), validation.ValidateData[requests.CreateCommentDTO](), commentsService.CreateComment)
                commentGroup.POST("/like/:commentID", validation.ValidatePathParams[requests.CommentIDPathParams](), commentsService.LikeComment)
                commentGroup.POST("/dislike/:commentID", validation.ValidatePathParams[requests.CommentIDPathParams](), commentsService.DislikeComment)
                commentGroup.PATCH("/current/:commentID", validation.ValidatePathParams[requests.CommentIDPathParams](), validation.ValidateData(requests.ValidateUpdateCommentDTO), commentsService.UpdateComment)
                commentGroup.DELETE("/current/:commentID", validation.ValidatePathParams[requests.CommentIDPathParams](), commentsService.DeleteCurrentUserComment)
                commentGroup.DELETE("/votes/current/:commentID", validation.ValidatePathParams[requests.CommentIDPathParams](), commentsService.RemoveCommentVote)

                adminOnly := commentGroup.Group("/admin", authSerivce.ValidateAdminUser)
                {
                    adminOnly.DELETE("/:commentID", validation.ValidatePathParams[requests.CommentIDPathParams](), commentsService.DeleteComment)
                }

            }
        }
    }

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
