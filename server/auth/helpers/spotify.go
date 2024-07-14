package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"github.com/Jack-Gitter/tunes/models/spotifyResponses"
)

func RetrieveInitialAccessToken(authorizationCode string) (*spotifyresponses.AccessTokenResponnse, error) {

    queryParamsMap := url.Values{}
    queryParamsMap.Add("grant_type", "authorization_code")
    queryParamsMap.Add("code", authorizationCode)
    queryParamsMap.Add("redirect_uri", os.Getenv("REDIRECT_URI"))
    queryParams := queryParamsMap.Encode()

    basicAuthToken := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
    encodedBasicAuthToken := base64.StdEncoding.EncodeToString([]byte(basicAuthToken))

    accessTokenRequest, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(queryParams)))

    if err != nil {
        return nil, err
    }

    accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    accessTokenRequest.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedBasicAuthToken))

    client := &http.Client{}
    resp, err := client.Do(accessTokenRequest) 

    if err != nil {
        return nil, err
    }

    accessTokenResponseBody := &spotifyresponses.AccessTokenResponnse{}
    json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

    return accessTokenResponseBody, nil
}

func RetrieveUserProfile(accessToken string) (*spotifyresponses.ProfileResponse, error) {

    nReq, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", &bytes.Buffer{})

    if err != nil {
        return nil, err
    }

    nReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    
    client := &http.Client{}
    nResp, err := client.Do(nReq)

    if err != nil {
        return nil, err
    }

    respJson2 := &spotifyresponses.ProfileResponse{}

    json.NewDecoder(nResp.Body).Decode(respJson2)

    return respJson2, nil

}

func RetreiveAccessTokenFromRefreshToken(spotifyRefreshToken string) (*spotifyresponses.RefreshTokenResponse, error) {

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
    resp, _ := client.Do(accessTokenRefreshRequest) 

    accessTokenResponseBody := &spotifyresponses.RefreshTokenResponse{}
    json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

    return accessTokenResponseBody, nil

}
