package jwt

import (
	"os"
	"time"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {}

type IJWTService interface {
    CreateAccessJWT(spotifyID string, username string, accessToken string, accessTokenExpiresAt int, role responses.Role) (string, error) 
    CreateRefreshJWT(spotifyRefreshToken string) (string, error) 
    ValidateAccessToken(accessTokenJWT string) (*jwt.Token, error) 
    ValidateRefreshToken(refreshTokenJWT string) (*jwt.Token, error) 
}

func(j *JWTService) CreateAccessJWT(spotifyID string, username string, accessToken string, accessTokenExpiresAt int, role responses.Role) (string, error) {

	claims := &requests.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tunes",
			Subject:   "bitch",
			Audience:  []string{"another bitch"},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        "garbage for now",
		},
		SpotifyID:            spotifyID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
		UserRole:             role,
		Username:             username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString, customerrors.WrapBasicError(err)
}

func(j *JWTService) CreateRefreshJWT(spotifyRefreshToken string) (string, error) {
	claims := &requests.RefreshJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tunes",
			Subject:   "bitch",
			Audience:  []string{"another bitch"},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        "garbage for now",
		},
		RefreshToken: spotifyRefreshToken,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return refreshString, customerrors.WrapBasicError(err)
}

func(j *JWTService) ValidateAccessToken(accessTokenJWT string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(accessTokenJWT, &requests.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return token, customerrors.WrapBasicError(err)
}

func(j *JWTService) ValidateRefreshToken(refreshTokenJWT string) (*jwt.Token, error) {
	token, e := jwt.ParseWithClaims(refreshTokenJWT, &requests.RefreshJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return token, customerrors.WrapBasicError(e)
}
