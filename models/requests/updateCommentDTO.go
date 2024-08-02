package requests

import (
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/gin-gonic/gin"
)

type UpdateCommentDTO struct {
    CommentText *string
}

func ValidateUpdateCommentDTO(ucdto UpdateCommentDTO, c *gin.Context) error {

    if ucdto.CommentText == nil {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "need to at least provide update comment text"}
    }

    return nil
    
}
