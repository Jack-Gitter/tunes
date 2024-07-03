package middlware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateUserJWT(c *gin.Context) {

    
    jwtCookie, err := c.Cookie("JWT")
    if err != nil {
        panic(err)
    }
    fmt.Println(jwtCookie)

    token, err := jwt.ParseWithClaims(jwtCookie, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, err)
    }
    
    userClaims := token.Claims.(*auth.JWTClaims)

    fmt.Println(userClaims)

}
