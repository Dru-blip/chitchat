package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func GlobalErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}

	code := http.StatusInternalServerError
	message := "Unable to process request"
	var details map[string]string

	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) {
		if tmp := sc.StatusCode(); tmp != 0 {
			code = tmp
		}
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		code = http.StatusUnprocessableEntity
		message = "Invalid Input"
		details = make(map[string]string)
		for _, e := range ve {
			field := strings.ToLower(e.Field())
			details[field] = fmt.Sprintf("%s is %s", field, e.Tag())
		}
	} else if he, ok := err.(*echo.HTTPError); ok {
		message = he.Message
	}

	obj := ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	var cErr error
	if c.Request().Method == http.MethodHead {
		cErr = c.NoContent(code)
	} else {
		cErr = c.JSON(code, obj)
	}

	if cErr != nil {
		c.Logger().Error(message)
	}
}
