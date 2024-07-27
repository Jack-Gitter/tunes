package requests

import (
	"net/http"

	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
)

type UpdatePostRequestDTO struct {
    Rating *int `binding:"gte=0,lte=5"`
	Text   *string  
}

func ValidateUpdatePostRequestDTO(req UpdatePostRequestDTO) error {
    if req.Rating == nil && req.Text == nil {
        return customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "one of the fields must be defined"}
    }
    if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 5) {
        return customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad range for rating"}
    }
    return nil

}
