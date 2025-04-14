package v1

import "github.com/labstack/echo/v4"

func Register(e *echo.Echo, h *Handler) {
	e.POST("/v1/book", h.SaveBook)
}
