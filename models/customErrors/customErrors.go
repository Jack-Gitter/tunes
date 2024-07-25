package customerrors

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)


type CustomError struct {
    StatusCode int
    Msg string
}

func (ce CustomError) Error() string {
    return ce.Msg
}

func ErrorHandlerMiddleware(c *gin.Context) {
    c.Next()
    if len(c.Errors) < 1 {
        return
    }
    firstError := c.Errors[0].Err

    if err, ok := firstError.(*CustomError); ok {
        c.JSON(err.StatusCode, err.Msg)
    } else {
        c.JSON(http.StatusInternalServerError, err.Error())
    }
}

func WrapBasicError(err error) error {

    if err == nil {
        return nil
    }

    customError := &CustomError{}
    customError.StatusCode = http.StatusInternalServerError
    customError.Msg = err.Error()

    if err == sql.ErrNoRows {
        customError.StatusCode = http.StatusNotFound
        customError.Msg = "Resource not found"
        return customError
    }

    if err, ok := err.(*pq.Error); ok {
        switch err.Code {
            case "23505": 
                customError.StatusCode = http.StatusConflict
                customError.Msg = "Duplicate resource cannot be created"
                return customError
            case "23503":
                customError.StatusCode = http.StatusNotFound
                customError.Msg = "DB key constraint violated"
                return customError
            default: 
                customError.StatusCode = http.StatusInternalServerError
                customError.Msg = err.Error()
                return customError
        }
    }


    if errors.Is(err, jwt.ErrTokenExpired) {
        customError.StatusCode = http.StatusUnauthorized
        customError.Msg = "Please refresh JWT"
    } else if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, jwt.ErrTokenUnverifiable) {
        customError.StatusCode = http.StatusForbidden
        customError.Msg = "JWT has been tampered with"
    }

    return customError

}

