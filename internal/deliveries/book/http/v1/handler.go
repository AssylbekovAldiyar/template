package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"libs/common/ctxconst"

	"template/internal/services/book"
	"template/pkg/deliveries"
	"template/pkg/reqresp"
)

type Handler struct {
	service book.Service
	timeout time.Duration
}

func NewV1(service book.Service, timeoutSeconds int) *Handler {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 60
	}

	return &Handler{
		service: service,
		timeout: time.Second * time.Duration(timeoutSeconds),
	}
}

func (h *Handler) SaveBook(c echo.Context) error {
	var request reqresp.SaveBookRequest
	if err := c.Bind(&request); err != nil {
		return deliveries.HandleEcho(c, err)
	}

	ctx, cancel := h.context(c)
	defer cancel()

	response, err := h.service.SaveBook(ctx, request)
	if err != nil {
		return deliveries.HandleEcho(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) context(c echo.Context) (context.Context, context.CancelFunc) {
	ctx := c.Request().Context()
	// todo: добавить мета инфу (реквест айди, пользователь и тд)
	ctx = ctxconst.SetRequestID(ctx, "test-request-id")
	ctx = ctxconst.SetUserID(ctx, "test-user-id")
	ctx = ctxconst.SetUserPhoneNumber(ctx, "test-phone-number")

	return context.WithTimeout(ctx, h.timeout)
}
