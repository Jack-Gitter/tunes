package posts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	//"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/gin-gonic/gin"
)


func CreatePostForCurrentUser(c *gin.Context) {

    spotifyID, spotifyIDExists := c.Get("spotifyID")
    spotifyAccessToken, spotifyAccessTokenExists := c.Get("spotifyAccessToken")

    post := &models.Post{}
    err := c.ShouldBindBodyWithJSON(post)

    if err != nil {
        panic(err)
    }

    if !spotifyIDExists || !spotifyAccessTokenExists {
        c.JSON(http.StatusUnauthorized, "user is not signed in (did i forget to pass the JWT in the middleware?)")
        return
    }


    url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", post.SongID)
    songRequest, _ := http.NewRequest(http.MethodGet, url, nil)
    songRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))

    client := &http.Client{}
    resp, err := client.Do(songRequest) 
    
    if resp.StatusCode != 200 {
        c.JSON(http.StatusBadRequest, "invalid spotify song ID")
        return
    }

    spotifySongResponse := &models.SongResponse{}
    bodyString, err := io.ReadAll(resp.Body)
    json.Unmarshal(bodyString, spotifySongResponse)

    if err != nil {
        c.JSON(http.StatusBadRequest, "invalid post request body for request 'post'")
        return
    }

    post.AlbumID = spotifySongResponse.Album.Id
    if len(spotifySongResponse.Album.Images) > 0 {
        post.AlbumArtURI = spotifySongResponse.Album.Images[0].Url
    }
    post.SongName = spotifySongResponse.Name
    post.AlbumName = spotifySongResponse.Album.Name

    err = db.CreatePost(post, spotifyID.(string))

    if err != nil {
        c.JSON(http.StatusInternalServerError, "failed to put post into db for some reason")
        return
    }

    c.JSON(http.StatusOK, post)

}
