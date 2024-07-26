package helpers

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/lib/pq"
)

func RunTransactionWithExponentialBackoff(transFunc func() error, retryTimes int) error {

    failureError := customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "Failed after retrying SQL statement"}
    backoff := 1.0

    for i := 0; i < retryTimes; i++ {

        fmt.Println("trying transaction!")
        err := transFunc()

        if err != nil {
            switch err := err.(type) {
                case *pq.Error:
                    if err.Code == "40001" {
                        fmt.Println("we are backing off!!!")
                        val := math.Pow(100, backoff)
                        backoff+=1
                        time.Sleep(time.Millisecond * time.Duration(val))
                        continue
                    } else {
                        return customerrors.WrapBasicError(err)
                    }
                case *customerrors.CustomError: 
                    return err
                default: 
                    return failureError
            }
        } else {
            return nil
        }

    }

        return failureError
}
