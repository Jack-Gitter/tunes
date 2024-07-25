package customerrors

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
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

    switch e := firstError.(type) {
        case CustomError: 
            c.JSON(e.StatusCode, e.Msg)
        default:
            c.JSON(http.StatusInternalServerError, e.Error())
    }
}

func WrapBasicError(err error) *CustomError {

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

    return customError

}

