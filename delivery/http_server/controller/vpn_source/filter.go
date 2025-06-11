package vpnsource

import (
	"net/http"

	"github.com/labstack/echo/v4"
	VPNSourceSvc "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	httperror "github.com/mohamadrezamomeni/momo/pkg/http_error"
	vpnSourceSerializer "github.com/mohamadrezamomeni/momo/serializer/vpn_source"
)

func (h *Handler) FilterVPNSources(c echo.Context) error {
	vpnSources, err := h.VPNSourceSvc.FilterVPNSources(&VPNSourceSvc.FilterVPNSourcesDto{})
	if err != nil {
		msg, code := httperror.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	vpnSourcesSerializer := make([]*vpnSourceSerializer.VPNSource, 0)

	for _, vpnSource := range vpnSources {
		vpnSourcesSerializer = append(vpnSourcesSerializer, &vpnSourceSerializer.VPNSource{
			Country: vpnSource.Country,
			English: vpnSource.English,
		})
	}
	return c.JSON(http.StatusOK, vpnSourceSerializer.FilterVPNSources{
		VPNSources: vpnSourcesSerializer,
	})
}
