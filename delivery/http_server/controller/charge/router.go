package charge

import (
	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/delivery/http_server/middleware"
)

func (h *Handler) SetRouter(v1 *echo.Group) {
	v1.POST("/charges/:id/approve", h.ApproveCharge, middleware.AccessCheck(h.authSvc, true))
}
