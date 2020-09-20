package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rghiorghisor/basic-go-rest-api/errors"
)

// appError represents the formatted error to be returned as the response body, in case this is needed.
type appError struct {
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// JSONAppErrorHandler is the middleware handling the overall error handling mechanism.
func JSONAppErrorHandler() gin.HandlerFunc {
	return handle(gin.ErrorTypeAny)
}

func handle(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		rawErrors := c.Errors

		if len(rawErrors) <= 0 {
			return
		}

		err := rawErrors[0].Err
		var parsedError *appError

		switch err.(type) {
		case *errors.Error:
			castError := err.(*errors.Error)
			parsedError = &appError{
				Code:    castError.Code,
				Message: castError.Message,
			}
		default:
			parsedError = &appError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
		}

		parsedError.Timestamp = time.Now()
		c.AbortWithStatusJSON(parsedError.Code, parsedError)

		return

	}
}
