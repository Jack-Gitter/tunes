package posts

import (
	"net/http"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/Jack-Gitter/tunes/server/posts/helpers"
	"github.com/gin-gonic/gin"
)

// need to prevent the user from posting the same song multiple times
func CreatePostForCurrentUser(c *gin.Context) {

    spotifyID, spotifyIDExists := c.Get("spotifyID")
    spotifyUsername, spotifyUsernameExists := c.Get("spotifyUsername")
    spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")

    if !spotifyIDExists || !spotifyAccessTokenExists || !spotifyUsernameExists {
        c.JSON(http.StatusUnauthorized, "user is not signed in (did i forget to pass the JWT in the middleware?)")
        return
    }

    post := &models.Post{}
    err := c.ShouldBindBodyWithJSON(post)

    if err != nil {
        c.JSON(http.StatusBadRequest, "bad post body bruh")
        return
    }

    hasPostedAlready, err := helpers.UserHasPostedSongAlready(spotifyID.(string), post.SongID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "bad")
        return
    }

    if hasPostedAlready {
        c.JSON(http.StatusBadRequest, "you have posted this shit already")
        return
    }

    spotifySongResponse, err := helpers.GetSongDetailsFromSpotify(post.SongID, spotifyAccessToken.(string))

    if err != nil {
        c.JSON(http.StatusBadRequest, "could not parse response from spotify API for get song")
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
        c.JSON(http.StatusInternalServerError, "failed to put post into db for some reason")
        return
    }

    c.JSON(http.StatusOK, post)

}
