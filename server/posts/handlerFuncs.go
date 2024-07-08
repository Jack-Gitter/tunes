package posts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
	"github.com/gin-gonic/gin"
)


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
        c.JSON(http.StatusBadRequest, "bad request body");
        return
    }

    url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", post.SongID)
    songRequest, err := http.NewRequest(http.MethodGet, url, nil)

    if err != nil {
        c.JSON(http.StatusInternalServerError, "bad http request to get a song from spotify")
        return
    }

    songRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))

    client := &http.Client{}
    resp, err := client.Do(songRequest) 
    
    if resp.StatusCode != 200 {
        respBodyBytes, _ := io.ReadAll(resp.Body)
        c.JSON(http.StatusBadRequest, string(respBodyBytes))
        return
    }

    spotifySongResponse := &models.SongResponse{}
    bodyString, err := io.ReadAll(resp.Body)
    json.Unmarshal(bodyString, spotifySongResponse)

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
