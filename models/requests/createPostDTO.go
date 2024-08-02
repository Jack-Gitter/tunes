package requests

import (
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/gin-gonic/gin"
)

type CreatePostDTO struct {
	SongID *string
	Rating *int
	Text   *string
}


func ValidateCreatePostDTO(createPostDTO CreatePostDTO, c *gin.Context) error {
    if createPostDTO.SongID == nil {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Must provide a spotifyID"}
    }
    if createPostDTO.Rating != nil && (*createPostDTO.Rating < 0 || *createPostDTO.Rating > 5) {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Rating must be between 0 and 5"}
    }

    return nil
}
