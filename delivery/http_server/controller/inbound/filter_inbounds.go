package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	inboundSerializer "github.com/mohamadrezamomeni/momo/serializer/inbound"
)

func (h *Handler) Filter(c echo.Context) error {
	var req inboundControllerDto.FilterInboundsDto
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.validation.ValidateFilteringInbounds(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	inbounds, err := h.inboundSvc.Filter(&inboundServiceDto.FilterInbounds{
		Domain:  req.Domain,
		Port:    req.Port,
		VPNType: entity.ConvertStringVPNTypeToEnum(req.VPNType),
		UserID:  req.UserID,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	res := &inboundSerializer.FilterInboundsSerializer{
		Inbounds: make([]*inboundSerializer.InboundSerializer, 0),
	}

	for _, inbound := range inbounds {
		res.Inbounds = append(res.Inbounds, &inboundSerializer.InboundSerializer{
			Protocol:     inbound.Protocol,
			Domain:       inbound.Domain,
			Port:         inbound.Port,
			UserID:       inbound.UserID,
			IsNotified:   inbound.IsNotified,
			IsBlock:      inbound.IsBlock,
			IsAssigned:   inbound.IsAssigned,
			Start:        inbound.Start,
			End:          inbound.End,
			Tag:          inbound.Tag,
			TrafficUsage: inbound.TrafficUsage,
			TrafficLimit: inbound.TrafficLimit,
			ChargeCount:  inbound.ChargeCount,
		})
	}
	return c.JSON(http.StatusOK, res)
}
