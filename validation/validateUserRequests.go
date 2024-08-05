package validation

import (
	"net/http"
	customerrors "github.com/Jack-Gitter/tunes/models/customErrors"
	"github.com/Jack-Gitter/tunes/models/dtos/requests"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/gin-gonic/gin"
)

func ValidateUserRequestDTO(req requests.UpdateUserRequestDTO, c *gin.Context) error {
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
            if roleToSet != responses.BASIC_USER {
                return false
            }
            return  true
        default: 
            return false
    }
}
