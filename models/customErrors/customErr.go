package customerrors

type CustomError struct {
    Code int
    E error
}

func (ce CustomError) Error() string {
    return ce.E.Error()
}


const (
    NEO_CONSTRAINT_ERROR = "Neo.ClientError.Schema.ConstraintValidationFailed" 
)
