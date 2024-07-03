package middlware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func testValidator() func(*jwt.Parser) {
    return func(parser *jwt.Parser) {
        
    }
}
func validateUserJWT(c *gin.Context) {

    token, err := jwt.ParseWithClaims("token", &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, "jwt has been tampered with!")
    }
    
    userClaims := token.Claims.(auth.JWTClaims)

    fmt.Println(userClaims)

}
