package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/oidiral/e-commerce/services/auth-svc/internal/response"
	"net/http"
)

func ApiErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

			var apiErr *response.ApiError
			if errors.As(err, &apiErr) {
				c.JSON(apiErr.Code, apiErr)
			} else {
				c.JSON(http.StatusInternalServerError, response.ApiError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
			}
			c.Abort()
		}
	}
}
