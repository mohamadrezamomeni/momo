package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundSerializer "github.com/mohamadrezamomeni/momo/serializer/inbound"
)

func (i *Handler) Filter(c echo.Context) error {
	inbounds, err := i.inboundSvc.Filter()
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
