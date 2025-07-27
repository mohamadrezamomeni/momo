package vpn

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) CreateVPN(c echo.Context) error {
	var req vpnControllerDto.CreateVPN
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.validation.ValidateCreatingVPN(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	_, err = h.vpnSvc.Create(&vpnServiceDto.CreateVPN{
		VpnType:   entity.ConvertStringVPNTypeToEnum(req.VpnType),
		UserCount: req.UserCount,
		Domain:    req.Domain,
		Port:      req.Port,
		Country:   req.Country,
		StartPort: req.StartPort,
		EndPort:   req.EndPort,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
