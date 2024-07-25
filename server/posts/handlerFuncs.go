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

func CreatePostForCurrentUser(c *gin.Context) {

    spotifyID, spotifyIDExists := c.Get("spotifyID")
    spotifyUsername, spotifyUsernameExists := c.Get("spotifyUsername")
    spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")


    if !spotifyIDExists || !spotifyAccessTokenExists || !spotifyUsernameExists {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
        return
    }

    createPostDTO := &requests.CreatePostDTO{}
    err := c.ShouldBindBodyWithJSON(createPostDTO)

    if createPostDTO.Rating < 0 || createPostDTO.Rating > 5 {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad body"})
        return
    }

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    spotifySongResponse, err := helpers.GetSongDetailsFromSpotify(createPostDTO.SongID, spotifyAccessToken.(string))

    if err != nil {
        c.AbortWithError(-1, err)
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
        c.AbortWithError(-1, err)
        return
    }


    c.JSON(http.StatusOK, resp)

}

func LikePost(c *gin.Context) {
    currentUserSpotifyID, found := c.Get("spotifyID")
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
        return
    }

    err := db.LikeOrDislikePost(currentUserSpotifyID.(string), spotifyID, songID, true)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }


    c.Status(http.StatusNoContent)
}

func DislikePost(c *gin.Context) {

    currentUserSpotifyID, found := c.Get("spotifyID")
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
        return
    }

    err := db.LikeOrDislikePost(currentUserSpotifyID.(string), spotifyID, songID, false)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.Status(http.StatusNoContent)
}


func GetAllPostsForUserByID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")
    createdAt := c.Query("createdAt")

    posts, err := getAllPosts(spotifyID, createdAt)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, posts)
}

func GetAllPostsForCurrentUser(c *gin.Context) {
    spotifyID, spotifyIDExists := c.Get("spotifyID")
    createdAt := c.Query("createdAt")

    if !spotifyIDExists {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt lookup"})
        return
    }

    posts, err := getAllPosts(spotifyID.(string), createdAt)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, posts)
}

func GetPostBySpotifyIDAndSongID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    post, err := db.GetUserPostByID(songID, spotifyID)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, post)
}

func GetPostCurrentUserBySongID(c *gin.Context) {

    currentUserSpotifyID, found := c.Get("spotifyID")
    songID := c.Param("songID")

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuckin jwt"})
        return
    }

    post, err := db.GetUserPostByID(songID, currentUserSpotifyID.(string))

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }



    c.JSON(http.StatusOK, post)
}

func DeletePostBySpotifyIDAndSongID(c *gin.Context) {

    requestorSpotifyID, found := c.Get("spotifyID")

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "fuckin jwt"})
        return
    }

    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    if requestorSpotifyID != spotifyID {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "canot do that"})
        return
    }

    err := db.DeletePost(songID, spotifyID)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.Status(http.StatusNoContent)

}


func DeletePostForCurrentUserBySongID(c *gin.Context) {

    
    requestorSpotifyID, found := c.Get("spotifyID")

    if !found {
        c.AbortWithError(-1, customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
        return
    }
    songID := c.Param("songID")

    err := db.DeletePost(songID, requestorSpotifyID.(string))

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }


    c.Status(http.StatusNoContent)

}

func UpdateCurrentUserPost(c *gin.Context) {

    spotifyID, exists := c.Get("spotifyID")
    spotifyUsername, uexists := c.Get("spotifyUsername")
    songID := c.Param("songID")
    updatePostReq := &requests.UpdatePostRequestDTO{}

    err := c.ShouldBindBodyWithJSON(updatePostReq)

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    if !exists || !uexists {
        c.AbortWithError(-1, err)
        return
    }

    if updatePostReq.Text == nil && updatePostReq.Rating == nil {
        c.AbortWithError(-1, err)
        return
    }

    preview, err := db.UpdatePost(spotifyID.(string), songID, updatePostReq.Text, updatePostReq.Rating, spotifyUsername.(string))

    if err != nil {
        c.AbortWithError(-1, err)
        return
    }

    c.JSON(http.StatusOK, preview)
}

func getAllPosts(spotifyID string, createdAt string) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {
    var t time.Time 
    if createdAt == "" {
        t = time.Now().UTC()
    } else {
        t, _ = time.Parse(time.UTC.String(), createdAt)
    }

    return db.GetUserPostsPreviewsByUserID(spotifyID, t)

}
