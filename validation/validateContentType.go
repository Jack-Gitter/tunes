package validation

import (
	"net/http"
	"github.com/Jack-Gitter/tunes/models/customerrors"
	"github.com/gin-gonic/gin"
)

func ValidateContentTypeJSON(c *gin.Context) {

    contentType := c.GetHeader("Content-Type")
    if contentType != "application/json" {
        c.Error(&customerrors.CustomError{StatusCode: http.StatusUnsupportedMediaType, Msg: "MIME type must be application/json"})
        c.Abort()
    }
    
}
