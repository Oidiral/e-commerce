package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/oidiral/e-commerce/services/auth-svc/internal/response"
)

func handleValidationError(ctx *gin.Context, err error) {
	var errs validator.ValidationErrors

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs = ve
	}

	if errs != nil {
		details := make(map[string]string, len(errs))
		for _, fe := range errs {
			var msg string
			switch fe.Tag() {
			case "required":
				msg = "обязательно для заполнения"
			case "email":
				msg = "должен быть корректным email"
			case "min":
				msg = "значение слишком короткое"
			default:
				msg = "некорректное значение"
			}
			details[fe.Field()] = msg
		}
		response.RespondWithError(ctx, http.StatusBadRequest,
			"ошибка валидации полей", details)
		return
	}
	response.RespondWithError(ctx, http.StatusBadRequest,
		"некорректный запрос", nil)
}
