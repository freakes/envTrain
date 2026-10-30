package server

import (
	"errors"
	"net/http"

	"api/internal/app/tasks/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func New(log *logrus.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.RequestID(), middleware.Recover(), middleware.BodyLimit("1M"))
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		status, message := http.StatusInternalServerError, "internal server error"

		var httpErr *echo.HTTPError
		switch {
		case errors.Is(err, domain.ErrNotFound):
			status, message = http.StatusNotFound, domain.ErrNotFound.Error()
		case errors.Is(err, domain.ErrInvalidTitle):
			status, message = http.StatusUnprocessableEntity, domain.ErrInvalidTitle.Error()
		case errors.As(err, &httpErr):
			status = httpErr.Code

			message = http.StatusText(status)
			if s, ok := httpErr.Message.(string); ok && status < 500 {
				message = s
			}
		}

		if status >= 500 {
			log.WithError(err).Error("request failed")
		}

		if writeErr := c.JSON(status, map[string]string{"error": message}); writeErr != nil {
			log.WithError(writeErr).Error("write response")
		}
	}
	e.GET("/hc", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })

	return e
}
