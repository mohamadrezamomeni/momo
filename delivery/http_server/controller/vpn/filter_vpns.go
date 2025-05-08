package vpn

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
	vpnSerializer "github.com/mohamadrezamomeni/momo/serializer/vpn"
)

func (h *Handler) Filter(c echo.Context) error {
	var req vpnControllerDto.FilterVPNs
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
		})
	}

	err := h.validation.ValidateFilterVPNs(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
		})
	}
	vpns, err := h.vpnSvc.Filter(&vpnServiceDto.FilterVPNs{
		Domain:  req.Domain,
		VPNType: entity.ConvertStringVPNTypeToEnum(req.VPNType),
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
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
		})
	}
	return c.JSON(http.StatusAccepted, filterVpnsSerializer)
}
