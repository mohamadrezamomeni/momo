package vpnsource

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.PUT("/vpn_sources/:country", h.Create, middleware.AccessCheck(h.authSvc, true))
	v1.GET("/vpn_sources", h.FilterVPNSources, middleware.AccessCheck(h.authSvc, true))
}
