package posts

import (
	"net/http"
	"time"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/Jack-Gitter/tunes/server/posts/helpers"
	"github.com/gin-gonic/gin"
)

// @Summary Creates a post for the current user
// @Description Creates a post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param CreatePostDTO body requests.CreatePostDTO true "Information required to create a post"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts [post]
// @Security Bearer
func CreatePostForCurrentUser(c *gin.Context) {

	spotifyID, spotifyIDExists := c.Get("spotifyID")
	spotifyUsername, spotifyUsernameExists := c.Get("spotifyUsername")
	spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")

	if !spotifyIDExists || !spotifyAccessTokenExists || !spotifyUsernameExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}

	createPostDTO := &requests.CreatePostDTO{}
	c.ShouldBindBodyWithJSON(createPostDTO)

	spotifySongResponse, err := helpers.GetSongDetailsFromSpotify(createPostDTO.SongID, spotifyAccessToken.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var albumImage string = ""
	if len(spotifySongResponse.Album.Images) > 0 {
		albumImage = spotifySongResponse.Album.Images[0].Url
	}

	resp, err := db.CreatePost(
		spotifyID.(string),
		createPostDTO.SongID,
		spotifySongResponse.Name,
		spotifySongResponse.Album.Id,
		spotifySongResponse.Album.Name,
		albumImage,
		createPostDTO.Rating,
		createPostDTO.Text,
		time.Now().UTC(),
		spotifyUsername.(string),
	)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	resp.Likes = []responses.UserIdentifer{}
	resp.Dislikes = []responses.UserIdentifer{}

	c.JSON(http.StatusOK, resp)

}

// @Summary Likes a post for the current user
// @Description Likes a post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "Song ID of the post to like"
// @Param songID path string true "Spotify ID of the user who posted the song"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/likes/{spotifyID}/{songID} [post]
// @Security Bearer
func LikePost(c *gin.Context) {
	currentUserSpotifyID, found := c.Get("spotifyID")
	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}

	err := db.LikeOrDislikePost(currentUserSpotifyID.(string), spotifyID, songID, true)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Dislikes a post for the current user
// @Description Dislikes a post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "Song ID of the post to dislike"
// @Param songID path string true "Spotify ID of the user who posted the song"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/dislikes/{spotifyID}/{songID} [post]
// @Security Bearer
func DislikePost(c *gin.Context) {

	currentUserSpotifyID, found := c.Get("spotifyID")
	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}

	err := db.LikeOrDislikePost(currentUserSpotifyID.(string), spotifyID, songID, false)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get all of a users post previews
// @Description Get all of a users post previews
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "The user whos posts are recieved. Value is a spotify ID"
// @Param createdAt query string false "Pagination Key. Format is UTC timestamp"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/previews/users/{spotifyID} [get]
// @Security Bearer
func GetAllPostsForUserByID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")
	createdAt := c.Query("createdAt")

	posts, err := getAllPosts(spotifyID, createdAt)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, posts)
}

// @Summary Get all of a users post previews
// @Description Get all of a users post previews
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "The user whos posts are recieved. Value is a spotify ID"
// @Param createdAt query string false "Pagination Key. Format is UTC timestamp"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/previews/users/current [get]
// @Security Bearer
func GetAllPostsForCurrentUser(c *gin.Context) {
	spotifyID, spotifyIDExists := c.Get("spotifyID")
	createdAt := c.Query("createdAt")

	if !spotifyIDExists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt lookup"})
		c.Abort()
		return
	}

	posts, err := getAllPosts(spotifyID.(string), createdAt)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, posts)
}

// @Summary Get apath specific post
// @Description Get a specific post
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "The user who posted the song"
// @Param songID path string true "The songID of the posted song"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/{spotifyID}/{songID} [get]
// @Security Bearer
func GetPostBySpotifyIDAndSongID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

	post, err := db.GetUserPostByID(songID, spotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Get a specific post for the current user
// @Description Get a specific post for the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/current/{songID} [get]
// @Security Bearer
func GetPostCurrentUserBySongID(c *gin.Context) {

	currentUserSpotifyID, found := c.Get("spotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuckin jwt"})
		c.Abort()
		return
	}

	post, err := db.GetUserPostByID(songID, currentUserSpotifyID.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Deletes a specific post. Only accessible to admins
// @Description Deletes a specific post. Only accessible to admins
// @Tags Posts
// @Accept json
// @Produce json
// @Param spotifyID path string true "The spotify ID of the user who posted the song"
// @Param songID query string true "The songID of the posted song"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/admin/{spotifyID}/{songID} [delete]
// @Security Bearer
func DeletePostBySpotifyIDAndSongID(c *gin.Context) {

	spotifyID := c.Param("spotifyID")
	songID := c.Param("songID")

	err := db.DeletePost(songID, spotifyID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Deletes a post made by the current user
// @Description Deletes a post made by the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/current/{songID} [delete]
// @Security Bearer
func DeletePostForCurrentUserBySongID(c *gin.Context) {

	requestorSpotifyID, found := c.Get("spotifyID")

	if !found {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
		c.Abort()
		return
	}
	songID := c.Param("songID")

	err := db.DeletePost(songID, requestorSpotifyID.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

// @Summary Updates a post made by the current user
// @Description Updates a post made by the current user
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Param UpdatePostDTO body requests.UpdatePostRequestDTO true "The fields to update"
// @Success 200 {object} responses.PostPreview
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/current/{songID} [patch]
// @Security Bearer
func UpdateCurrentUserPost(c *gin.Context) {

	spotifyID, exists := c.Get("spotifyID")
	spotifyUsername, uexists := c.Get("spotifyUsername")
	songID := c.Param("songID")
	updatePostReq := &requests.UpdatePostRequestDTO{}

	c.ShouldBindBodyWithJSON(updatePostReq)

	if !exists || !uexists {
		c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "no jwt"})
		c.Abort()
		return
	}

	preview, err := db.UpdatePost(spotifyID.(string), songID, updatePostReq, spotifyUsername.(string))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, preview)
}

// @Summary Removes a vote for the current user on a post
// @Description Removes a vote for the current user on a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Param posterSpotifyID path string true "The user who posted the post spotify ID"
// @Success 204
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/votes/current/{posterSpotifyID}/{songID} [delete]
// @Security Bearer
func RemovePostVote(c *gin.Context) {
	voterSpotifyID, found := c.Get("spotifyID")
	posterSpotifyID := c.Param("posterSpotifyID")
	songID := c.Param("songID")

	if !found {
		c.Error(customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "forgot to set JWT"})
	}

	err := db.RemoveVote(voterSpotifyID.(string), posterSpotifyID, songID)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)

}

func getAllPosts(spotifyID string, createdAt string) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {

	var t time.Time = time.Now().UTC()
    var err error

	if createdAt != "" {
        t, err = time.Parse(time.RFC3339, createdAt)

        if err != nil {
            return nil, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid time format"}
        }

	}

	return db.GetUserPostsPreviewsByUserID(spotifyID, t)

}

// @Summary Gets the comments of a post
// @Description Gets the comments of a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param songID path string true "The songID of the posted song"
// @Param spotifyID path string true "The user who posted the post spotify ID"
// @Param createdAt query string false "Pagination Key. In the form of UTC timestamp"
// @Success 200 {object} responses.PaginationResponse[[]responses.Comment, time.Time]
// @Failure 400 {string} string 
// @Failure 404 {string} string 
// @Failure 500 {string} string 
// @Router /posts/current/{songID} [delete]
// @Security Bearer
func GetPostCommentsPaginated(c *gin.Context) {
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")
    createdAt := c.Query("createdAt")

    var t time.Time
    var err error

	if createdAt != "" {
		t, err = time.Parse(time.RFC3339, createdAt)

        if err != nil {
            c.Error(customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "invalid time format"})
            c.Abort()
            return
        }
	}

    resp, err := db.GetPostCommentsPaginated(spotifyID, songID, t)

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.JSON(http.StatusFound, resp)
}
