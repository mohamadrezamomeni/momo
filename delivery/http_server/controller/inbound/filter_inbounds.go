package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	inboundSerializer "github.com/mohamadrezamomeni/momo/serializer/inbound"
)

func (h *Handler) Filter(c echo.Context) error {
	var req inboundControllerDto.FilterInboundsDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
		})
	}

	err := h.validation.ValidateFilteringInbounds(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
		})
	}
	inbounds, err := h.inboundSvc.Filter(&inboundServiceDto.FilterInbounds{
		Domain:  req.Domain,
		Port:    req.Port,
		VPNType: entity.ConvertStringVPNTypeToEnum(req.VPNType),
		UserID:  req.UserID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "something went wrong",
		})
	}
	res := &inboundSerializer.FilterInboundsSerializer{
		Inbounds: make([]*inboundSerializer.InboundSerializer, 0),
	}

	for _, inbound := range inbounds {
		res.Inbounds = append(res.Inbounds, &inboundSerializer.InboundSerializer{
			Protocol:   inbound.Protocol,
			Domain:     inbound.Domain,
			Port:       inbound.Port,
			UserID:     inbound.UserID,
			IsNotified: inbound.IsNotified,
			IsBlock:    inbound.IsBlock,
			IsAssigned: inbound.IsAssigned,
			Start:      inbound.Start,
			End:        inbound.End,
			Tag:        inbound.Tag,
		})
	}
	return c.JSON(http.StatusOK, res)
}
