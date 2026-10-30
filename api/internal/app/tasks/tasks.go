// Package tasks is the feature composition root, matching the source project's Init pattern.
package tasks

import (
	delivery "api/internal/app/tasks/delivery/http/echo"
	"api/internal/app/tasks/domain"
	"api/internal/app/tasks/usecase"
	"github.com/labstack/echo/v4"
)

func Init(api *echo.Group, repo domain.Repository) {
	delivery.New(api.Group("/tasks"), usecase.New(repo)).Expose()
}
