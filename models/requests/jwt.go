package requests

import (
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    SpotifyID string
    AccessToken string
    RefreshToken string
    AccessTokenExpiresAt int
    UserRole responses.Role
    Username string
    jwt.RegisteredClaims
}
