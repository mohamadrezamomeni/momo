package host

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hostDto "github.com/mohamadrezamomeni/momo/dto/controller/host"
	"github.com/mohamadrezamomeni/momo/dto/service/host"
	"github.com/mohamadrezamomeni/momo/entity"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	hostSerializer "github.com/mohamadrezamomeni/momo/serializer/host"
	hostTransformer "github.com/mohamadrezamomeni/momo/transformer/host"
)

func (h *Handler) FilterHosts(c echo.Context) error {
	var req hostDto.FilterHostsDto
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	hostStatuses, err := hostTransformer.TransformStringsToHostStatus(req.Statuses)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	hosts, err := h.hostSvc.Filter(&host.FilterHosts{
		Status: hostStatuses,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	resp := &hostSerializer.FilterHosts{
		Hosts: make([]*hostSerializer.Host, 0),
	}
	for _, host := range hosts {
		resp.Hosts = append(resp.Hosts, &hostSerializer.Host{
			Domain: host.Domain,
			Port:   host.Port,
			Rank:   host.Rank,
			Status: entity.HostStatusString(host.Status),
		})
	}

	return c.JSON(http.StatusAccepted, resp)
}
