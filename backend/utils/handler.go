package utils

import (
	"log/slog"

	"github.com/labstack/echo/v5"
)

// LogRequest logs the start of a request with basic info.
func LogRequest(logger *slog.Logger, c *echo.Context) {
	logger.Info("request started",
		slog.String("method", c.Request().Method),
		slog.String("path", c.Request().URL.Path),
	)
}

// HandleError logs an error and returns a JSON response with the error message and status code.
func HandleError(logger *slog.Logger, c *echo.Context, err error, code int) error {
	logger.Error(err.Error(),
		slog.Int("status", code),
		slog.String("path", c.Request().URL.Path),
		slog.String("method", c.Request().Method),
	)
	return c.JSON(code, map[string]string{
		"error": err.Error(),
	})
}
