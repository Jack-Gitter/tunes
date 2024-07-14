package requests

import "github.com/Jack-Gitter/tunes/models/responses"

type UpdateUserRequestDTO struct {
    Bio string
    Role responses.Role
}
