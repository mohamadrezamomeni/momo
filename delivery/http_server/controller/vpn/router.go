package vpn

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.POST("/vpns", h.CreateVPN, middleware.AccessCheck(h.authSvc, true))
	v1.GET("/vpns", h.Filter, middleware.AccessCheck(h.authSvc, true))
	v1.PATCH("/vpns/:id", h.Update, middleware.AccessCheck(h.authSvc, true))
}
