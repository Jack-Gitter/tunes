package db

import (
	"database/sql"
	"net/http"

	"github.com/lib/pq"
)

type DBError struct {
    StatusCode int
    Err error
}

// returns error if it is 500, nil if not

func (dbe DBError) Error() string {
    return dbe.Err.Error()
}

func HandleDatabaseError(err error) error {
    dbError := &DBError{}
    dbError.Err = err

    if err == sql.ErrNoRows {
        dbError.StatusCode = http.StatusNotFound
        return dbError
    }

    if err, ok := err.(*pq.Error); ok {
        switch err.Code {
            case "10000": 
                dbError.StatusCode = http.StatusBadRequest
                return dbError
        }
    }
    panic("we should never be here!")
}

