package posts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/Jack-Gitter/tunes/server/posts/helpers"
	"github.com/gin-gonic/gin"
)

func CreatePostForCurrentUser(c *gin.Context) {

    spotifyID, spotifyIDExists := c.Get("spotifyID")
    spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")

    if !spotifyIDExists || !spotifyAccessTokenExists {
        c.JSON(http.StatusUnauthorized, "No JWT data found for the current user")
        return
    }

    createPostDTO := &requests.CreatePostDTO{}
    err := c.ShouldBindBodyWithJSON(createPostDTO)

    if createPostDTO.Rating < 0 || createPostDTO.Rating > 5 {
        c.JSON(http.StatusBadRequest, "please rate from 0-5!")
        return
    }

    if err != nil {
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }

    hasPostedAlready, err := helpers.UserHasPostedSongAlready(spotifyID.(string), createPostDTO.SongID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if hasPostedAlready {
        c.JSON(http.StatusBadRequest, "post with songID is already found for user")
        return
    }

    spotifySongResponse, err := helpers.GetSongDetailsFromSpotify(createPostDTO.SongID, spotifyAccessToken.(string))

    if err != nil {
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }


    var albumImage string = ""
    if len(spotifySongResponse.Album.Images) > 0 {
        albumImage = spotifySongResponse.Album.Images[0].Url
    }

    post, err := db.CreatePost(
        spotifyID.(string),
        createPostDTO.SongID,
        spotifySongResponse.Name, 
        spotifySongResponse.Album.Id, 
        spotifySongResponse.Album.Name,
        albumImage,
        createPostDTO.Rating,
        createPostDTO.Text,
        time.Now().UTC(),
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, post)

}

func LikePost(c *gin.Context) {

    /*currentUserSpotifyID, found := c.Get("spotifyID")
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    if !found {
        c.JSON(http.StatusInternalServerError, "nope to jwt")
        return
    }

    postPreview, found, err := db.LikePostForUser(currentUserSpotifyID, spotifyID, songID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusBadRequest, "no post found")
        return
    }

    c.JSON(http.StatusOK, postPreview)*/

}

func GetAllPostsForUserByID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")
    createdAt := c.Query("createdAt")

    posts, err := getAllPosts(spotifyID, createdAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, posts)
}

func GetAllPostsForCurrentUser(c *gin.Context) {
    spotifyID, spotifyIDExists := c.Get("spotifyID")
    createdAt := c.Query("createdAt")

    if !spotifyIDExists {
        c.JSON(http.StatusUnauthorized, "No JWT data found for the current user")
        return
    }

    posts, err := getAllPosts(spotifyID.(string), createdAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, posts)
}

func GetPostBySpotifyIDAndSongID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    post, found, err := db.GetUserPostByID(songID, spotifyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "could not find post with that userid and songid in the database")
        return
    }


    c.JSON(http.StatusOK, post)
}

func GetPostCurrentUserBySongID(c *gin.Context) {

    currentUserSpotifyID, found := c.Get("spotifyID")
    songID := c.Param("songID")

    if !found {
        c.JSON(http.StatusInternalServerError, "spotifyID not set in JWT middleware")
        return
    }

    post, found, err := db.GetUserPostByID(songID, currentUserSpotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "could not find post with that userid and songid in the database")
        return
    }


    c.JSON(http.StatusOK, post)
}

func DeletePostBySpotifyIDAndSongID(c *gin.Context) {

    requestorSpotifyID, found := c.Get("spotifyID")

    if !found {
        c.JSON(http.StatusInternalServerError, "no spotify ID found for user making request (did I forget to pass it in the middleware?)")
    }
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    if requestorSpotifyID != spotifyID {
        c.JSON(http.StatusBadRequest, "cannot delete a post that is not your own! (unless you're admin, tbd)")
        return
    }

    _, found, err := db.DeletePost(songID, spotifyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "something went wrong with deletion")
        return
    }

    if !found {
        c.JSON(http.StatusBadRequest, "post for that user has not been found!")
        return
    }

    c.JSON(http.StatusOK, "post deleted")


}


func DeletePostForCurrentUserBySongID(c *gin.Context) {

    
    requestorSpotifyID, found := c.Get("spotifyID")

    if !found {
        c.JSON(http.StatusInternalServerError, "no spotify ID found for user making request (did I forget to pass it in the middleware?)")
    }
    songID := c.Param("songID")

    _, found, err := db.DeletePost(songID, requestorSpotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, "something went wrong with deletion")
        return
    }

    if !found {
        c.JSON(http.StatusBadRequest, "post for that user has not been found!")
        return
    }

    c.JSON(http.StatusOK, "post deleted")

}

func UpdateCurrentUserPost(c *gin.Context) {

    spotifyID, exists := c.Get("spotifyID")
    songID := c.Param("songID")
    updatePostReq := &requests.UpdatePostRequestDTO{}

    err := c.ShouldBindBodyWithJSON(updatePostReq)

    if err != nil {
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "bad json body")
        return
    }

    if !exists {
        c.JSON(http.StatusBadRequest, "need jwt")
        return
    }

    preview, found, err := db.UpdatePost(spotifyID.(string), songID, updatePostReq.Text, updatePostReq.Rating)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusBadRequest, "could not find post or user")
        return
    }

    c.JSON(http.StatusOK, preview)
}

func getAllPosts(spotifyID string, createdAt string) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {
    var t time.Time 
    if createdAt == "" {
        t = time.Now().UTC()
    } else {
        t, _ = time.Parse(time.RFC3339, createdAt)
    }

    posts, err := db.GetUserPostsPreviewsByUserID(spotifyID, t)

    if err != nil {
        return nil, err
    }

    return posts, nil

}
