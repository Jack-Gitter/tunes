package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
    SpotifyID string
    AccessToken string
    RefreshToken string
    AccessTokenExpiresAt int
    UserRole string
    jwt.RegisteredClaims
}
