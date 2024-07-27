package requests

import (
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/responses"
)

// the body fields must be pointers, because the zero value for pointers is nill. We will be able
// to properly determine whether or not users have requested to update or change a resource
type UpdateUserRequestDTO struct {
	Bio  *string
	Role *responses.Role 
}

func ValidateUserRequestDTO(req UpdateUserRequestDTO) error {
    if req.Bio == nil && req.Role == nil {
        return customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad body"}
    }
    if req.Role != nil && ! responses.IsValidRole(*req.Role) {
        return customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Invalid role"}
    }
    return nil
}
