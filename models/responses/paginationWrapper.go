package responses

type PaginationResponse[T any, U any] struct {
    PaginationKey U 
    DataResponse T
}
