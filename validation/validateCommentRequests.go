package validation

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/gin-gonic/gin"
)

func ValidateUpdateCommentDTO(ucdto requests.UpdateCommentDTO, c *gin.Context) error {

    if ucdto.CommentText == nil {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "need to at least provide update comment text"}
    }

    return nil
    
}
