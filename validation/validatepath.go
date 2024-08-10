package validation

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/gin-gonic/gin"
)
func ValidatePathParams[T any]() func(c *gin.Context) {

    return func(c *gin.Context) {
        val := new(T)

        err := c.ShouldBindUri(&val)

        if err != nil {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad path parameter value"})
            c.Abort()
            return
        }
        c.Next()
    }

}
