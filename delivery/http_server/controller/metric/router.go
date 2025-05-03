package metric

import "github.com/labstack/echo/v4"

func (h *Handler) SetRouter(r *echo.Group) {
	r.GET("/v1/healthz", h.healthCheck)
}
