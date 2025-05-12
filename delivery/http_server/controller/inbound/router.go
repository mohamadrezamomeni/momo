package inbound

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.POST("/inbounds", h.CreateInbound, middleware.AccessCheck(h.authSvc, true))
	v1.GET("/inbounds", h.Filter, middleware.AccessCheck(h.authSvc, true))
	v1.POST("/inbounds/:id/block", h.Block, middleware.AccessCheck(h.authSvc, true))
	v1.POST("/inbounds/:id/unblock", h.UnBlock, middleware.AccessCheck(h.authSvc, true))
	v1.POST("/inbounds/:id/extend", h.ExtendInbound, middleware.AccessCheck(h.authSvc, true))
	v1.PATCH("/inbounds/:id", h.UpdateInbound, middleware.AccessCheck(h.authSvc, true))
}
