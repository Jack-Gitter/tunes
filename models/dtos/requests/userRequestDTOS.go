package requests

import "github.com/Jack-Gitter/tunes/models/dtos/responses"

type UpdateUserRequestDTO struct {
	Bio  *string
    Email *string
	UserRole *responses.Role 
}
