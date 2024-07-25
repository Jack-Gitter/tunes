package db

import (
	"database/sql"
	"net/http"

	"github.com/lib/pq"
)

type DBError struct {
    StatusCode int
    Msg string
}

// returns error if it is 500, nil if not

func (dbe DBError) Error() string {
    return dbe.Msg
}

func HandleDatabaseError(err error) error {
    dbError := &DBError{}
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
    panic("we should never be here!")
}
