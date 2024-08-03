package requests

import (
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	SpotifyID            string
	AccessToken          string
	AccessTokenExpiresAt int
	UserRole             responses.Role
	Username             string
	jwt.RegisteredClaims
}

type RefreshJWTClaims struct {
	RefreshToken string
	jwt.RegisteredClaims
}

type RefreshJWTDTO struct {
	AccessToken  string
	RefreshToken string
}
