package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/spotifyResponses"
)

func GetSongDetailsFromSpotify(songID string, spotifyAccessToken string) (*spotifyresponses.SongResponse, error) {

	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", songID)
	songRequest, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	songRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))

	client := &http.Client{}
	resp, err := client.Do(songRequest)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Invalid spotify song information"}
		} else {
			return nil, customerrors.WrapBasicError(errors.New("spotify request for a song failed without 200"))
		}
	}

	spotifySongResponse := &spotifyresponses.SongResponse{}
	bodyString, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyString, spotifySongResponse)

	return spotifySongResponse, nil

}
