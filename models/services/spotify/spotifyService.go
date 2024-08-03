package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
)

type SpotifyService struct { }

type ISpotifyService interface {
    RetrieveInitialAccessToken(authorizationCode string) (*responses.AccessTokenResponnse, error) 
    RetrieveUserProfile(accessToken string) (*responses.ProfileResponse, error) 
    RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken string) (*responses.RefreshTokenResponse, error) 
}

func(s *SpotifyService) RetrieveInitialAccessToken(authorizationCode string) (*responses.AccessTokenResponnse, error) {

	queryParamsMap := url.Values{}
	queryParamsMap.Add("grant_type", "authorization_code")
	queryParamsMap.Add("code", authorizationCode)
	queryParamsMap.Add("redirect_uri", os.Getenv("REDIRECT_URI"))
	queryParams := queryParamsMap.Encode()

	basicAuthToken := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	encodedBasicAuthToken := base64.StdEncoding.EncodeToString([]byte(basicAuthToken))

	accessTokenRequest, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(queryParams)))

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	accessTokenRequest.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedBasicAuthToken))

	client := &http.Client{}
	resp, err := client.Do(accessTokenRequest)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	defer resp.Body.Close()
	accessTokenResponseBody := &responses.AccessTokenResponnse{}
	json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

	return accessTokenResponseBody, nil
}

func(s *SpotifyService) RetrieveUserProfile(accessToken string) (*responses.ProfileResponse, error) {

	nReq, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", &bytes.Buffer{})

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	nReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	nResp, err := client.Do(nReq)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	defer nResp.Body.Close()
	respJson2 := &responses.ProfileResponse{}

	json.NewDecoder(nResp.Body).Decode(respJson2)

	return respJson2, nil

}

func(s *SpotifyService) RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken string) (*responses.RefreshTokenResponse, error) {

	queryParamsMap := url.Values{}
	queryParamsMap.Add("grant_type", "refresh_token")
	queryParamsMap.Add("refresh_token", spotifyRefreshToken)
	queryParams := queryParamsMap.Encode()

	basicAuthToken := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	encodedBasicAuthToken := base64.StdEncoding.EncodeToString([]byte(basicAuthToken))

	accessTokenRefreshRequest, _ := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(queryParams)))
	accessTokenRefreshRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	accessTokenRefreshRequest.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedBasicAuthToken))

	client := &http.Client{}
	resp, err := client.Do(accessTokenRefreshRequest)

	if err != nil {
		return nil, customerrors.WrapBasicError(err)
	}

	defer resp.Body.Close()

	accessTokenResponseBody := &responses.RefreshTokenResponse{}
	json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

	return accessTokenResponseBody, nil

}
