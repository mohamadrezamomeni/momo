package vpn

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vpnControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/vpn"
	vpnServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn"
	"github.com/mohamadrezamomeni/momo/entity"
)

func (h *Handler) CreateVPN(c echo.Context) error {
	var req vpnControllerDto.CreateVPN
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input is not valid",
		})
	}

	err := h.validation.ValidateCreatingVPN(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
		})
	}

	_, err = h.vpnSvc.Create(&vpnServiceDto.CreateVPN{
		VpnType:   entity.ConvertStringVPNTypeToEnum(req.VpnType),
		UserCount: req.UserCount,
		Domain:    req.Domain,
		Port:      req.Port,
	})
	return c.NoContent(http.StatusNoContent)
}
