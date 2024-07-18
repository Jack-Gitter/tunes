package responses

import "time"

type PaginationResponse[T any] struct {
    PaginationKey time.Time
    More bool
    DataResponse T
}
