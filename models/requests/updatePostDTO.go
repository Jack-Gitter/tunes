package requests

import (
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/gin-gonic/gin"
)

type UpdatePostRequestDTO struct {
    Rating *int 
	Review   *string  
}

func ValidateUpdatePostRequestDTO(req UpdatePostRequestDTO, c *gin.Context) error {
    if req.Rating == nil && req.Review == nil {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "one of the fields must be defined"}
    }
    if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 5) {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad range for rating"}
    }
    return nil

}
