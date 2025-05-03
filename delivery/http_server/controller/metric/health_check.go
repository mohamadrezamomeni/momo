package metric

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) healthCheck(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
