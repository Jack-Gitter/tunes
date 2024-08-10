package validation

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/gin-gonic/gin"
)

func ValidateUpdatePostRequestDTO(req requests.UpdatePostRequestDTO, c *gin.Context) error {
    if req.Rating == nil && req.Review == nil {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "one of the fields must be defined"}
    }
    if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 5) {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad range for rating"}
    }
    return nil

}

func ValidateCreatePostDTO(createPostDTO requests.CreatePostDTO, c *gin.Context) error {
    if createPostDTO.SongID == nil {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Must provide a spotifyID"}
    }
    if createPostDTO.Rating != nil && (*createPostDTO.Rating < 0 || *createPostDTO.Rating > 5) {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Rating must be between 0 and 5"}
    }

    return nil
}
