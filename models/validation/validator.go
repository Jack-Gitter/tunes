package validation

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/gin-gonic/gin"
)

func ValidateData[T any](funcs ...func(T) error) func(c *gin.Context) {

    return func(c *gin.Context) {
        validationObject := new(T)
        err := c.ShouldBindBodyWithJSON(&validationObject)
        if err != nil {
            c.Error(&customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: err.Error()})
            c.Abort()
            return
        }

        for _, function := range funcs {
            err := function(*validationObject)
            if err != nil {
                c.Error(err)
                c.Abort()
                return
            }
        }
        c.Next()
    }

}
