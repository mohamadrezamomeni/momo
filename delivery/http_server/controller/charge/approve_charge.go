package charge

import (
	"net/http"

	"github.com/labstack/echo/v4"
	identifyChargeDto "github.com/mohamadrezamomeni/momo/dto/controller/charge"
	httperror "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) ApproveCharge(c echo.Context) error {
	var req identifyChargeDto.ApproveChargeDto
	if err := c.Bind(&req); err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.chargeSvc.ApproveCharge(req.ID)
	if err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	return c.NoContent(http.StatusOK)
}
