package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Jack-Gitter/tunes/customerrors"
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models"
)


func GetSongDetailsFromSpotify(songID string, spotifyAccessToken string) (*models.SongResponse, error) {


    url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", songID)
    songRequest, err := http.NewRequest(http.MethodGet, url, nil)

    if err != nil {
        return nil, err
    }

    songRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))

    client := &http.Client{}
    resp, err := client.Do(songRequest) 
    
    if resp.StatusCode != 200 {
        return nil, errors.New("spotify request for a song failed without 200")
    }

    spotifySongResponse := &models.SongResponse{}
    bodyString, err := io.ReadAll(resp.Body)
    json.Unmarshal(bodyString, spotifySongResponse)

    return spotifySongResponse, nil

}

func UserHasPostedSongAlready(spotifyID string, songID string) (bool, error) {

    _, err := db.GetUserPostById(songID, spotifyID)

    if err != nil {
        if customError, ok := err.(customerrors.TunesError); ok && customError.ErrorType == customerrors.NoDatabaseRecordsFoundError {
            return false, nil
        } 
        return false, err
    }

    return true, nil
}
