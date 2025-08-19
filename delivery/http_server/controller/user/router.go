package user

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.GET("/users", h.filterUsers, middleware.AccessCheck(h.authSvc, true))
	v1.POST("/users/:id/approve", h.ApproveUser, middleware.AccessCheck(h.authSvc, true))
}
