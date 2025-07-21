package vpnpackage

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.POST("/vpn_packages", h.Create, middleware.AccessCheck(h.authSvc, true))
	v1.GET("/vpn_packages", h.FilterPackages, middleware.AccessCheck(h.authSvc, true))
}
