package contoller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/service"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func HandleError(c *gin.Context, err error) {
	var resp ErrorResponse
	var status int

	switch {
	case errors.Is(err, service.ErrNotFound):
		resp = ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Resource not found",
		}
		status = http.StatusNotFound
	case errors.Is(err, service.ErrBadRequest):
		resp = ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Bad request",
		}
		status = http.StatusBadRequest
	default:
		resp = ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		}
		status = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(status, resp)
}
