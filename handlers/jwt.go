package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

)

func RetrieveAccessToken(authorizationCode string) *AccessTokenResponnse {

    queryParamsMap := url.Values{}
    queryParamsMap.Add("grant_type", "authorization_code")
    queryParamsMap.Add("code", authorizationCode)
    queryParamsMap.Add("redirect_uri", os.Getenv("REDIRECT_URI"))
    queryParams := queryParamsMap.Encode()

    basicAuthToken := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
    encodedBasicAuthToken := base64.StdEncoding.EncodeToString([]byte(basicAuthToken))

    accessTokenRequest, _ := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(queryParams)))
    accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    accessTokenRequest.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedBasicAuthToken))

    client := &http.Client{}
    resp, _ := client.Do(accessTokenRequest) 

    accessTokenResponseBody := &AccessTokenResponnse{}
    json.NewDecoder(resp.Body).Decode(accessTokenResponseBody)

    return accessTokenResponseBody
}

func RetrieveUserProfile(accessToken string) *ProfileResponse {

    nReq, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", &bytes.Buffer{})
    nReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    
    client := &http.Client{}
    nResp, _ := client.Do(nReq)

    respJson2 := &ProfileResponse{}

    json.NewDecoder(nResp.Body).Decode(respJson2)

    return respJson2

}

