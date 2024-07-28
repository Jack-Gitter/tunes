package requests

import (
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
)

type CreatePostDTO struct {
	SongID string
	Rating int
	Text   string
}


func ValidateCreatePostDTO(createPostDTO CreatePostDTO) error {
    if createPostDTO.Rating < 0 || createPostDTO.Rating > 5 {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad body"}
    }
    return nil
}
