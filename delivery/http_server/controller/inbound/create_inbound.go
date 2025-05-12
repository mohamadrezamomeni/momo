package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	"github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	inboundSerializer "github.com/mohamadrezamomeni/momo/serializer/inbound"
	timeTransformer "github.com/mohamadrezamomeni/momo/transformer/time"
)

func (h *Handler) CreateInbound(c echo.Context) error {
	var req inboundControllerDto.CreateInbound

	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.validation.ValidateCreatingInbound(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	startTime, _ := timeTransformer.ConvertStrToTime(req.Start)
	endTime, _ := timeTransformer.ConvertStrToTime(req.End)
	inboundCreated, err := h.inboundSvc.Create(&inbound.CreateInbound{
		UserID:       req.UserID,
		Start:        startTime,
		TrafficLimit: req.TrafficLimit,
		End:          endTime,
		VPNType:      entity.ConvertStringVPNTypeToEnum(req.VPNType),
		ServerType:   entity.High,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	return c.JSON(http.StatusAccepted, &inboundSerializer.CredateInboundSerializer{
		ID: inboundCreated.ID,
	})
}
