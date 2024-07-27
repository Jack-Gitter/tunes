package requests

import "github.com/Jack-Gitter/tunes/models/responses"

// the body fields must be pointers, because the zero value for pointers is nill. We will be able
// to properly determine whether or not users have requested to update or change a resource
type UpdateUserRequestDTO struct {
	Bio  *string
	Role *responses.Role 
}
