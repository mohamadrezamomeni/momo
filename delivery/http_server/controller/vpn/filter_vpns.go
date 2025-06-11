package vpn

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	vpnSerializer "github.com/mohamadrezamomeni/momo/serializer/vpn"
)

func (h *Handler) Filter(c echo.Context) error {
	var req vpnControllerDto.FilterVPNs
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.validation.ValidateFilterVPNs(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	vpns, err := h.vpnSvc.Filter(&vpnServiceDto.FilterVPNs{
		Domain:  req.Domain,
		VPNType: entity.ConvertStringVPNTypeToEnum(req.VPNType),
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	filterVpnsSerializer := &vpnSerializer.FilterVpnsSerializer{
		VPNs: make([]*vpnSerializer.VPNSerializer, 0),
	}
	for _, vpn := range vpns {
		filterVpnsSerializer.VPNs = append(filterVpnsSerializer.VPNs, &vpnSerializer.VPNSerializer{
			ID:        vpn.ID,
			ApiPort:   vpn.ApiPort,
			UserCount: vpn.UserCount,
			Domain:    vpn.Domain,
			VPNType:   entity.VPNTypeString(vpn.VPNType),
			Country:   vpn.Country,
		})
	}
	return c.JSON(http.StatusAccepted, filterVpnsSerializer)
}
