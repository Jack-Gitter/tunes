package models

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    SpotifyID string
    AccessToken string
    RefreshToken string
    AccessTokenExpiresAt int
    UserRole string
    jwt.RegisteredClaims
}

func (c JWTClaims) Validate() error {
    if c.SpotifyID != "blah" {
        // here, check if the access token is expeired for the user
        return errors.New("if you are seeing this, we still need to implement custom JWT claim validation")
    }
    return nil
}

func createNewJWT() (JWTClaims, error) {

    return JWTClaims{}, nil
}
