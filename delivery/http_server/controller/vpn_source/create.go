package vpnsource

import (
	"net/http"

	"github.com/labstack/echo/v4"
	VPNSourceDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn_source"
	VPNSourceSvc "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	httperror "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) Create(c echo.Context) error {
	var req VPNSourceDto.CreateVPNSourceDto
	if err := c.Bind(&req); err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	if err := h.vpnSourceValidator.ValidateUpsert(req); err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.VPNSourceSvc.Create(&VPNSourceSvc.CreateVPNSourceDto{
		Country: req.Country,
		English: req.English,
	})
	if err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
