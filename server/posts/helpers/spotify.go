package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/spotifyResponses"
)


func GetSongDetailsFromSpotify(songID string, spotifyAccessToken string) (*spotifyresponses.SongResponse, error) {

    url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", songID)
    songRequest, err := http.NewRequest(http.MethodGet, url, nil)

    if err != nil {
        return nil, err
    }

    songRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))

    client := &http.Client{}
    resp, err := client.Do(songRequest) 
    
    if resp.StatusCode != 200 {
        bodyBytes, _ := io.ReadAll(resp.Body)
        fmt.Println(string(bodyBytes))
        return nil, errors.New("spotify request for a song failed without 200")
    }

    spotifySongResponse := &spotifyresponses.SongResponse{}
    bodyString, err := io.ReadAll(resp.Body)
    json.Unmarshal(bodyString, spotifySongResponse)

    return spotifySongResponse, nil

}

func UserHasPostedSongAlready(spotifyID string, songID string) (bool, error) {

    _, foundPost, err := db.GetUserPostByID(songID, spotifyID)

    if err != nil {
        return false, err
    }

    if !foundPost {
        return false, nil
    }

    return true, nil
}
