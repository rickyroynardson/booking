package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rickyroynardson/booking/internal/entity"
	"github.com/rickyroynardson/booking/internal/service"
)

type ShowHandler struct {
	service   *service.ShowService
	validator *validator.Validate
}

func NewShowHandler(service *service.ShowService, validator *validator.Validate) *ShowHandler {
	return &ShowHandler{
		service,
		validator,
	}
}

func (h *ShowHandler) Create(c echo.Context) error {
	var body entity.CreateShowRequest
	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := h.service.Create(c.Request().Context(), body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, "success create show")
}
