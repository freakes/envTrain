package delivery

import (
	"encoding/json"
	"io"
	"net/http"

	"api/internal/app/tasks/domain"
	"github.com/labstack/echo/v4"
)

type delivery struct {
	group *echo.Group
	uc    domain.Usecase
}

func New(group *echo.Group, uc domain.Usecase) domain.Delivery {
	return &delivery{group: group, uc: uc}
}
func (d *delivery) Expose() {
	d.group.POST("", d.Create)
	d.group.GET("", d.List)
	d.group.GET("/:id", d.Get)
}
func (d *delivery) Create(c echo.Context) error {
	var req domain.CreateRequest

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON body")
	}

	if err := decoder.Decode(new(any)); err != io.EOF {
		return echo.NewHTTPError(http.StatusBadRequest, "expected one JSON object")
	}

	task, err := d.uc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Location", "/api/tasks/"+task.ID)

	return c.JSON(http.StatusCreated, task)
}
func (d *delivery) Get(c echo.Context) error {
	task, err := d.uc.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}
func (d *delivery) List(c echo.Context) error {
	tasks, err := d.uc.List(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tasks)
}
