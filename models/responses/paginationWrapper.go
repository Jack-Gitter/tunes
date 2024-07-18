package responses

import "time"

type PaginationResponse[T any] struct {
    PaginationKey time.Time
    DataResponse T
}
