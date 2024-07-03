package helpers

import (
	"os"
	"time"

	"github.com/Jack-Gitter/tunes/models"
	"github.com/golang-jwt/jwt/v5"
)


func CreateAccessJWT(spotifyID string, accessToken string, refreshToken string, accessTokenExpiresAt int) (string, error) {

    claims := &models.JWTClaims{
        RegisteredClaims: jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
        },
        SpotifyID: spotifyID,
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        AccessTokenExpiresAt: accessTokenExpiresAt,
        UserRole: "user",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

    return tokenString, err
}

func CreateRefreshJWT() (string, error) {

    claimsForRefresh := &jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsForRefresh)
    refreshString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

    return refreshString, err
}

func ValidateAccessToken(accessTokenJWT string) (*jwt.Token, error) {
    token, err := jwt.ParseWithClaims(accessTokenJWT, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    return token, err
}

func ValidateRefreshToken(refreshTokenJWT string) (*jwt.Token, error) {
    token, e := jwt.Parse(refreshTokenJWT, func (token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    return token, e
}
