package usertier

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.POST("/tiers/:tier/users/:user_id/assign", h.Create, middleware.AccessCheck(h.authSvc, true))
	v1.DELETE("/tiers/:tier/users/:user_id", h.Delete, middleware.AccessCheck(h.authSvc, true))
	v1.GET("/users/:user_id/tiers", h.FilterTiersByUser, middleware.AccessCheck(h.authSvc, true))
}
