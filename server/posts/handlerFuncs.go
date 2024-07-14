package posts

import (
	"net/http"

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
