package posts

import (
	"net/http"
	"time"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/Jack-Gitter/tunes/server/posts/helpers"
	"github.com/gin-gonic/gin"
)

func CreatePostForCurrentUser(c *gin.Context) {

    spotifyID, spotifyIDExists := c.Get("spotifyID")
    spotifyUsername, spotifyUsernameExists := c.Get("spotifyUsername")
    spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")

    if !spotifyIDExists || !spotifyAccessTokenExists || !spotifyUsernameExists {
        c.JSON(http.StatusUnauthorized, "No JWT data found for the current user")
        return
    }

    post := &models.Post{}
    err := c.ShouldBindBodyWithJSON(post)

    if err != nil {
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }

    hasPostedAlready, err := helpers.UserHasPostedSongAlready(spotifyID.(string), post.SongID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    if hasPostedAlready {
        c.JSON(http.StatusBadRequest, "post with songID is already found for user")
        return
    }

    spotifySongResponse, err := helpers.GetSongDetailsFromSpotify(post.SongID, spotifyAccessToken.(string))

    if err != nil {
        c.JSON(http.StatusBadRequest, err.Error())
        return
    }

    post.AlbumID = spotifySongResponse.Album.Id
    post.SongName = spotifySongResponse.Name
    post.AlbumName = spotifySongResponse.Album.Name
    post.SpotifyID = spotifyID.(string)
    post.Username = spotifyUsername.(string)
    post.Timestamp = time.Now().UTC()

    if len(spotifySongResponse.Album.Images) > 0 {
        post.AlbumArtURI = spotifySongResponse.Album.Images[0].Url
    }

    err = db.CreatePost(post, spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, post)

}

func GetPostBySpotifyIDAndSongID(c *gin.Context) {

    spotifyID := c.Param("spotifyID")
    songID := c.Param("songID")

    post, found, err := db.GetUserPostById(songID, spotifyID)

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

    post, found, err := db.GetUserPostById(songID, currentUserSpotifyID.(string))

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


func DeletePostForCurrentUserBySongID (c *gin.Context) {

    
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
