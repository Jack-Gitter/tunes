package db

import (
	"math"
	"net/http"
	"time"

	"github.com/Jack-Gitter/tunes/models/customErrors"
)

func RunTransactionWithExponentialBackoff(transFunc func() error, retryTimes int) error {

    failureError := &customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Failed after retrying SQL statement"}
    backoff := 1.0

    for i := 0; i < retryTimes; i++ {

        err := transFunc()

        if err != nil {
            switch err := err.(type) {
                case *customerrors.CustomError: 
                    if err.StatusCode == 40001 {
                        val := math.Pow(100, backoff)
                        backoff+=1
                        time.Sleep(time.Millisecond * time.Duration(val))
                        continue
                    } else {
                        return err
                    }
                default: 
                    return failureError
            }
        } else {
            return nil
        }

    }

        return failureError
}
