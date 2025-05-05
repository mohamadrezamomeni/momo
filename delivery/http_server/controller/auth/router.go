package auth

import "github.com/labstack/echo/v4"

func (h *Handler) SetRouter(api *echo.Group) {
	api.POST("/auth/login", h.Login)
}
