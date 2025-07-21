package vpnpackage

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vpnPackageControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn_package"
	vpnPackageServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_package"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	vpnPackageSerializer "github.com/mohamadrezamomeni/momo/serializer/vpn_package"
)

func (h *Handler) Create(c echo.Context) error {
	var req vpnPackageControllerDto.CreateVPNPackage
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	err := h.vpnPackageValidator.CreateVPNPackage(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	vpnPakcage, err := h.vpnPackageSvc.Create(
		&vpnPackageServiceDto.CreateVPNPackage{
			Days:              req.Days,
			Months:            req.Months,
			Price:             req.Price,
			IsActive:          req.IsActive,
			TrafficLimitTitle: req.TrafficLimitTitle,
			TrafficLimit:      req.TrafficLimit,
			PriceTitle:        req.PriceTitle,
		},
	)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	return c.JSON(http.StatusCreated, &vpnPackageSerializer.CreateVPNPackage{
		ID: vpnPakcage.ID,
	})
}
