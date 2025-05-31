package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func RespondWithError(c *gin.Context, code int, message string, details interface{}) {
	apiErr := &ApiError{
		Code:    code,
		Message: message,
		Details: details,
	}
	c.Error(apiErr)
	c.Abort()
}

func BadRequest(c *gin.Context, message string, details interface{}) {
	RespondWithError(c, http.StatusBadRequest, message, details)
}
