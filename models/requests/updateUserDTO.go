package requests

import (
	"net/http"

	"github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/responses"
	"github.com/gin-gonic/gin"
)

// the body fields must be pointers, because the zero value for pointers is nill. We will be able
// to properly determine whether or not users have requested to update or change a resource
type UpdateUserRequestDTO struct {
	Bio  *string
	UserRole *responses.Role 
}

func ValidateUserRequestDTO(req UpdateUserRequestDTO, c *gin.Context) error {
    userRole, found := c.Get("userRole")

    if !found {
        return &customerrors.CustomError{StatusCode: http.StatusInternalServerError, Msg: "jwt mess"}
    }
    if req.Bio == nil && req.UserRole == nil {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "bad body"}
    }
    if req.UserRole != nil && !responses.IsValidRole(*req.UserRole) {
        return &customerrors.CustomError{StatusCode: http.StatusBadRequest, Msg: "Invalid role"}
    }
    if req.UserRole != nil && !CanSetRole(userRole.(responses.Role), *req.UserRole) {
        return &customerrors.CustomError{StatusCode: http.StatusForbidden, Msg: "Cannot upgrade your own role"}
    }
    return nil
}

func CanSetRole(currentUserRole responses.Role, roleToSet responses.Role) bool {
    switch currentUserRole{
        case responses.ADMIN:
            return true
        case responses.MODERATOR:
            if roleToSet == responses.ADMIN {
                return false
            }
            return true
        case responses.BASIC_USER:
            return  false
        default: 
            return false
    }
}
