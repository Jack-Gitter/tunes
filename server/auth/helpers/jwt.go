package helpers

import (
	"os"
	"time"

	"github.com/Jack-Gitter/tunes/models/requests"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/golang-jwt/jwt/v5"
)


func CreateAccessJWT(spotifyID string, username string, accessToken string, refreshToken string, accessTokenExpiresAt int, role responses.Role) (string, error) {

    claims := &requests.JWTClaims{
        RegisteredClaims: jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           //ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
           ExpiresAt: &jwt.NumericDate{Time: time.Now()},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
        },
        SpotifyID: spotifyID,
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        AccessTokenExpiresAt: accessTokenExpiresAt,
        UserRole: role,
        Username: username,
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
    token, err := jwt.ParseWithClaims(accessTokenJWT, &requests.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
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
