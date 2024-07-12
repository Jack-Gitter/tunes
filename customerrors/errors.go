package customerrors

const (
    SpotifyRequestError = iota
    Neo4jDatabaseRequestError 
    NoDatabaseRecordsFoundError
)

type TunesError struct {
    ErrorType int
    Err error
} 

func (e TunesError) Error() string {
    return e.Err.Error()
}

func (e TunesError) GetErrorType() int {
    return e.ErrorType
}

