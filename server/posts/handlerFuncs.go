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
    spotifyUsername, spotifyUsernameExists := c.Get("spotifyUsername")
    spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")


    if !spotifyIDExists || !spotifyAccessTokenExists || !spotifyUsernameExists {
        c.JSON(http.StatusInternalServerError, "No JWT data found for the current user")
        return
    }

    createPostDTO := &requests.CreatePostDTO{}
    err := c.ShouldBindBodyWithJSON(createPostDTO)

    if createPostDTO.Rating < 0 || createPostDTO.Rating > 5 {
        c.JSON(http.StatusBadRequest, "Please provide a rating 0-5")
        return
    }

    if err != nil {
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }

    spotifySongResponse, err := helpers.GetSongDetailsFromSpotify(createPostDTO.SongID, spotifyAccessToken.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    var albumImage string = ""
    if len(spotifySongResponse.Album.Images) > 0 {
        albumImage = spotifySongResponse.Album.Images[0].Url
    }

    resp, collision, err := db.CreatePost(
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
        fmt.Println(err.Error())
        c.JSON(http.StatusInternalServerError, "being lazy just search for this error")
        return
    }

    if collision {
        c.JSON(http.StatusBadRequest, "cant post same song twice")
        return
    }

    c.JSON(http.StatusOK, resp)

}

func LikePost(c *gin.Context) {
    currentUserSpotifyID, found := c.Get("spotifyID")
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    if !found {
        c.JSON(http.StatusInternalServerError, "Did not set spotifyID in JWT middleware")
        return
    }

    found, err := db.LikePostForUser(currentUserSpotifyID.(string), spotifyID, songID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "Could not find post or user with specified values")
        return
    }

    c.Status(http.StatusNoContent)
}

func GetAllPostsForUserByID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")
    createdAt := c.Query("createdAt")

    posts, foundUser, err := getAllPosts(spotifyID, createdAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !foundUser {
        c.JSON(http.StatusNotFound, "Could not find user with ID")
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

    posts, foundUser, err := getAllPosts(spotifyID.(string), createdAt)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !foundUser {
        c.JSON(http.StatusNotFound, "Could not find user with ID")
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
        c.JSON(http.StatusNotFound, "Could not find post with corresponding userid and songid")
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
        c.JSON(http.StatusNotFound, "Could not find post with that userid and songid")
        return
    }


    c.JSON(http.StatusOK, post)
}

func DeletePostBySpotifyIDAndSongID(c *gin.Context) {

    requestorSpotifyID, found := c.Get("spotifyID")

    if !found {
        c.JSON(http.StatusInternalServerError, "No spotifyID set from JWT middleware")
    }
    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    if requestorSpotifyID != spotifyID {
        c.JSON(http.StatusForbidden, "Cannot delete post that is not your own, unless you are an admin")
        return
    }

    found, err := db.DeletePost(songID, spotifyID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "Post for user with songID could not be found")
        return
    }

    c.Status(http.StatusNoContent)


}


func DeletePostForCurrentUserBySongID(c *gin.Context) {

    
    requestorSpotifyID, found := c.Get("spotifyID")

    if !found {
        c.JSON(http.StatusInternalServerError, "No spotifyID variable set in JWT middleware")
    }
    songID := c.Param("songID")

    found, err := db.DeletePost(songID, requestorSpotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, "Post with specified songID for user could not be found")
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
        fmt.Println(err.Error())
        c.JSON(http.StatusBadRequest, "bad json body")
        return
    }

    if !exists || !uexists {
        c.JSON(http.StatusBadRequest, "need jwt")
        return
    }

    preview, found, err := db.UpdatePost(spotifyID.(string), songID, updatePostReq.Text, updatePostReq.Rating, spotifyUsername.(string))

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

func getAllPosts(spotifyID string, createdAt string) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], bool, error) {
    var t time.Time 
    if createdAt == "" {
        t = time.Now().UTC()
    } else {
        t, _ = time.Parse(time.RFC3339, createdAt)
    }

    return db.GetUserPostsPreviewsByUserID(spotifyID, t)

}
