package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
)

// AccessLogger retrieves a new middleware that logs the GIN requests.
func AccessLogger() gin.HandlerFunc {
	return newLogger(*logger.Access)
}

func newLogger(logger logger.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		stop := time.Now()
		latency := stop.Sub(start)
		statusCode := c.Writer.Status()

		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		msg1 := fmt.Sprintf("%3d | %13v | %8v | %-7s %#v %s",
			statusCode,
			latency,
			dataLength,
			c.Request.Method,
			path,
			c.Errors.ByType(gin.ErrorTypePrivate).String())

		if statusCode >= http.StatusInternalServerError {
			logger.Errore(msg1)
		} else if statusCode >= http.StatusBadRequest {
			logger.Warn(msg1)
		} else {
			logger.Info(msg1)
		}
	}
}
