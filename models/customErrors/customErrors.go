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
	Msg        string
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
	switch err := firstError.(type) {
	case *CustomError:
		c.JSON(err.StatusCode, err.Msg)
	default:
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

    if ok := wrapJWTErrors(err, customError); ok {
        return customError
    }
    if ok := wrapSQLErrors(err, customError); ok {
        return customError
    }
    if ok := wrapPostgresDriverErrors(err, customError); ok {
        return customError
    }

	return customError

}

func wrapPostgresDriverErrors(err error, customError *CustomError) bool {  
	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505":
			customError.StatusCode = http.StatusConflict
			customError.Msg = "Duplicate resource cannot be created"
            return true
		case "23503":
			customError.StatusCode = http.StatusNotFound
			customError.Msg = "DB key constraint violated"
            return true
        case "40001": 
            customError.StatusCode = 40001
            customError.Msg = "retry transaction!"
            return true
		}
	}
    return false
}

func wrapSQLErrors(err error, customError *CustomError) bool {
	if errors.Is(err, sql.ErrNoRows) {
		customError.StatusCode = http.StatusNotFound
		customError.Msg = "Resource not found"
		return true
    }
    return false

}

func wrapJWTErrors(err error, customError *CustomError) bool {
	if errors.Is(err, jwt.ErrTokenExpired) {
		customError.StatusCode = http.StatusUnauthorized
		customError.Msg = "Please refresh JWT"
        return true
	} else if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, jwt.ErrTokenUnverifiable) {
		customError.StatusCode = http.StatusForbidden
		customError.Msg = "JWT has been tampered with"
        return true
	}
    return false
}
