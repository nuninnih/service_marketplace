package common

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ValidationErrors(err error) interface{} {
	var messages []string

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				messages = append(messages,
					fmt.Sprintf("%s is required", e.Field()))
			case "email":
				messages = append(messages,
					fmt.Sprintf("%s must be a valid email", e.Field()))
			case "max":
				messages = append(messages,
					fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param()))
			case "gt":
				messages = append(messages,
					fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param()))
			case "lt":
				messages = append(messages,
					fmt.Sprintf("%s must be less than %s", e.Field(), e.Param()))
			default:
				messages = append(messages,
					"Invalid Specification")
			}
		}
	}

	return messages
}

func SuccessResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"data":    data,
	}
}

func CompleteSuccessResponse(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(
		statusCode,
		map[string]interface{}{
			"message": http.StatusText(statusCode),
			"data":    data,
		},
	)
}

func CompleteErrorResponse(c echo.Context, statusCode int, message interface{}) error {
	return c.JSON(
		statusCode,
		map[string]interface{}{
			"message": message,
		},
	)
}
