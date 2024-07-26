package db

import (
	"math"
	"net/http"
	"time"

	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/lib/pq"
)

func RunTransactionWithExponentialBackoff(transFunc func() error, retryTimes int) error {

    backoff := 1.0
    for i := 0; i < retryTimes; i++ {

        err := transFunc()

        if err != nil {
            if err, ok := err.(*pq.Error); ok {
                if err.Code == "40001" {
                    val := math.Pow(100, backoff)
                    time.Sleep(time.Millisecond * time.Duration(val))
                    continue
                } else {
                    return err
                }
            }
        } else {
            return nil
        }
    }

    return customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Failed after retrying SQL statement"}
}
