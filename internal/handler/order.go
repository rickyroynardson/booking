package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rickyroynardson/booking/internal/entity"
	"github.com/rickyroynardson/booking/internal/service"
)

type OrderHandler struct {
	service   *service.OrderService
	validator *validator.Validate
}

func NewOrderHandler(service *service.OrderService, validator *validator.Validate) *OrderHandler {
	return &OrderHandler{
		service,
		validator,
	}
}

func (h *OrderHandler) Book(c echo.Context) error {
	var req entity.BookOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := h.service.Book(c.Request().Context(), req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "book success")
}
