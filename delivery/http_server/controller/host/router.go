package host

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.POST("/hosts", h.CreateHost, middleware.AccessCheck(h.authSvc, true))
	v1.GET("/hosts", h.FilterHosts, middleware.AccessCheck(h.authSvc, true))
}
