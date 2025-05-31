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

func RespondWithError(ctx *gin.Context, code int, message string, details interface{}) {
	apiErr := &ApiError{
		Code:    code,
		Message: message,
		Details: details,
	}
	ctx.Error(apiErr)
	ctx.Abort()
}

func BadRequest(ctx *gin.Context, message string, details interface{}) {
	RespondWithError(ctx, http.StatusBadRequest, message, details)
}
