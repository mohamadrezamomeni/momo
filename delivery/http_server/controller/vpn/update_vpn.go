package vpn

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (v *Handler) Update(c echo.Context) error {
	var req vpnControllerDto.UpdateVPN
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := v.validation.ValidateUpdatingVPN(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	vpnStatus := entity.ConvertVPNStatusLabelToVPNStatus(
		req.VPNStatusLabel,
	)
	err = v.vpnSvc.Update(req.ID, &vpnServiceDto.Update{
		VPNStatus: vpnStatus,
		ApiPort:   req.ApiPort,
		Domain:    req.Domain,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
