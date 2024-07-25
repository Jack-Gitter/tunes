package customerrors

import (
//	"database/sql"
	"errors"
//	"net/http"

	"github.com/gin-gonic/gin"
//	"github.com/lib/pq"
)


var PermissionDeniedError = errors.New("Cannot perform this operation with your user role status")
var InternalServerError = errors.New("The server has encountered an error")
var BadJSONBodyError = errors.New("Your JSON request is either malformed, or it has invalid data in it")


// returns error if it is 500, nil if not


/*func HandleDatabaseError(err error) *CustomError {
    dbError := &CustomError{}
    dbError.StatusCode = http.StatusBadRequest
    dbError.Msg = err.Error()

    if err == sql.ErrNoRows {
        dbError.StatusCode = http.StatusNotFound
        dbError.Msg = "Resource not found"
        return dbError
    }

    if err, ok := err.(*pq.Error); ok {
        switch err.Code {
            case "23505": 
                dbError.StatusCode = http.StatusConflict
                dbError.Msg = "Duplicate resource cannot be created"
                return dbError
            case "23503":
                dbError.StatusCode = http.StatusNotFound
                dbError.Msg = "DB key constraint violated"
                return dbError
            default: 
                dbError.StatusCode = http.StatusInternalServerError
                dbError.Msg = err.Error()
                return dbError
        }
    }

    return dbError
}*/

func ErrorHandlerMiddleware (c *gin.Context) {
    c.Next()
    if len(c.Errors) < 1 {
        return
    }
    /*firstError := c.Errors[0].Err
    if err, ok := firstError.(*pq.Error); ok {
        respCode, resp := handlePQError(err)
        c.JSON(respCode, resp)

    }*/
    // if it is not a special type of error, than we have set the status code already. Just reurn the error message
    // set the response header and the message, and return
}

