package handler

import (
	"errors"
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

func (h *ShowHandler) FindAll(c echo.Context) error {
	var req entity.FindAllShowRequest

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 15
	}
	if req.Search != "" {
		req.Search = "%" + req.Search + "%"
	}

	res, err := h.service.FindAll(c.Request().Context(), req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *ShowHandler) FindById(c echo.Context) error {
	id := c.Param("id")

	res, err := h.service.FindById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrShowNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
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
