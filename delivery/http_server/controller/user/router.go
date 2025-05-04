package user

import "github.com/labstack/echo/v4"

func (h *Handler) SetHandler(v1 *echo.Group) {
	v1.POST("/users", h.Create)
}
