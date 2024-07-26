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
        c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad jwt"})
        c.Abort()
        return
    }

    createPostDTO := &requests.CreatePostDTO{}
    err := c.ShouldBindBodyWithJSON(createPostDTO)

    if err != nil {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Bad JSON Body"})
        c.Abort()
        return
    }

    if createPostDTO.Rating < 0 || createPostDTO.Rating > 5 {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad body"})
        c.Abort()
        return
    }

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

func UpdateCurrentUserPost(c *gin.Context) {

    spotifyID, exists := c.Get("spotifyID")
    spotifyUsername, uexists := c.Get("spotifyUsername")
    songID := c.Param("songID")
    updatePostReq := &requests.UpdatePostRequestDTO{}

    err := c.ShouldBindBodyWithJSON(updatePostReq)

    if err != nil {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Bad JSON Body"})
        c.Abort()
        return
    }

    if !exists || !uexists {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "bad body"})
        c.Abort()
        return
    }

    if updatePostReq.Text == nil && updatePostReq.Rating == nil {
        c.Error(err)
        c.Abort()
        return
    }

    if (updatePostReq.Rating != nil) && *updatePostReq.Rating > 5 || *updatePostReq.Rating < 0 {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Value must be between 0 and 5 for raitng"})
        c.Abort()
        return
    }

    preview, err := db.UpdatePost(spotifyID.(string), songID, updatePostReq.Text, updatePostReq.Rating, spotifyUsername.(string))

    if err != nil {
        c.Error(err)
        c.Abort()
        return
    }

    c.JSON(http.StatusOK, preview)
}

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
    var t time.Time 
    if createdAt == "" {
        t = time.Now().UTC()
    } else {
        t, _ = time.Parse(time.RFC3339, createdAt)
    }

    return db.GetUserPostsPreviewsByUserID(spotifyID, t)

}

