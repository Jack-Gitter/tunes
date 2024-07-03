package middlware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/Jack-Gitter/tunes/server/auth/spotifyHelpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateUserJWT(c *gin.Context) {
    
    jwtCookie, err := c.Cookie("JWT")
    if err != nil {
        panic(err)
    }

    token, err := jwt.ParseWithClaims(jwtCookie, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        spotifyID := token.Claims.(*auth.JWTClaims).SpotifyID
        fmt.Println(spotifyID)
        // send a request to refreshJWT which refreshes the JWT and sends it back to the user. Send the spotify ID there
        c.JSON(http.StatusBadRequest, err.Error())
    }
    
    //userClaims := token.Claims.(*auth.JWTClaims)


}

func refreshJWT(c *gin.Context) {

    refreshToken, err := c.Cookie("REFRESH_JWT")
    spotifyID := c.Query("spotifyID")

    if err != nil {
        panic(err)
    }

    _, e := jwt.Parse(refreshToken, func (token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if e != nil {
        // reroute the user to login screen, because refresh token has expired
        c.JSON(http.StatusUnauthorized, "the refresh token has expired. Please log out and log back in again")
    }

    accessTokenResponse := spotifyHelpers.RetrieveAccessToken(spotifyID)
    userProfileResponse := spotifyHelpers.RetrieveUserProfile(accessTokenResponse.Access_token)

    claims := &auth.JWTClaims{
        RegisteredClaims: jwt.RegisteredClaims{
           Issuer: "tunes", 
           Subject: "bitch",
           Audience: []string{"another bitch"},
           ExpiresAt: &jwt.NumericDate{Time: time.Now()},
           NotBefore: &jwt.NumericDate{Time: time.Now()},
           IssuedAt: &jwt.NumericDate{Time: time.Now()},
           ID: "garbage for now",
        },
        SpotifyID: userProfileResponse.Id,
        AccessToken: accessTokenResponse.Access_token,
        RefreshToken: accessTokenResponse.Refresh_token,
        AccessTokenExpiresAt: accessTokenResponse.Expires_in,
        UserRole: "user",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

    c.SetCookie("JWT", tokenString, 3600, "/", "localhost", false, true)
    c.Status(http.StatusOK)
}
    // get the query parameter of the 
    // here we have to get the spotifyID, and then make a request to our database to get the user object, and re-populate the JWT
    // how do we get the spotifyID though if the normal JWT has expired?

}
