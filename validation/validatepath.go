package validation

import (
	"fmt"
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/gin-gonic/gin"
)
func ValidatePathParams[T any]() func(c *gin.Context) {

    return func(c *gin.Context) {
        val := new(T)

        err := c.ShouldBindUri(&val)

        if err != nil {
            fmt.Println(err.Error())
            c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad path parameter value"})
            c.Abort()
            return
        }
        c.Next()
    }

}
