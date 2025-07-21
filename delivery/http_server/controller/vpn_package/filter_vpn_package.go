package vpnpackage

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vpnPackageControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn_package"
	vpnpackage "github.com/mohamadrezamomeni/momo/dto/service/vpn_package"
	"github.com/mohamadrezamomeni/momo/entity"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	vpnPackageSerializer "github.com/mohamadrezamomeni/momo/serializer/vpn_package"
)

func (h *Handler) FilterPackages(c echo.Context) error {
	var req vpnPackageControllerDto.FilterVPNPackages
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	vpnPackages, err := h.vpnPackageSvc.Filter(&vpnpackage.FilterVPNPackage{})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	serializer := &vpnPackageSerializer.FilterVPNPackagesSerializer{
		VPNPackages: h.makeVPNPackages(vpnPackages),
	}
	return c.JSON(http.StatusOK, serializer)
}

func (h *Handler) makeVPNPackages(vpnPackages []*entity.VPNPackage) []*vpnPackageSerializer.VPNPackageSerializer {
	vpnPackagesSerializer := make([]*vpnPackageSerializer.VPNPackageSerializer, 0)
	for _, vpnPackage := range vpnPackages {
		vpnPackagesSerializer = append(vpnPackagesSerializer, &vpnPackageSerializer.VPNPackageSerializer{
			Price:             vpnPackage.Price,
			PriceTitle:        vpnPackage.PriceTitle,
			TrafficLimit:      vpnPackage.TrafficLimit,
			TrafficLimitTitle: vpnPackage.TrafficLimitTitle,
			IsActive:          vpnPackage.IsActive,
			Months:            vpnPackage.Months,
			Days:              vpnPackage.Days,
			ID:                vpnPackage.ID,
		})
	}
	return vpnPackagesSerializer
}
