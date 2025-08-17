package tier

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.PUT("/tiers/:name", h.Create, middleware.AccessCheck(h.authSvc, true))
	v1.PATCH("/tiers/:name", h.Update, middleware.AccessCheck(h.authSvc, true))
	v1.GET("/tiers", h.Filter, middleware.AccessCheck(h.authSvc, true))
}
