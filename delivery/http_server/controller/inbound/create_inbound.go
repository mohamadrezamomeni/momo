package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	"github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	inboundSerializer "github.com/mohamadrezamomeni/momo/serializer/inbound"
	timeTransformer "github.com/mohamadrezamomeni/momo/transformer/time"
)

func (h *Handler) CreateInbound(c echo.Context) error {
	var req inboundControllerDto.CreateInbound

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input is Wrong",
		})
	}

	err := h.validation.ValidateCreatingInbound(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input is Wrong",
		})
	}

	startTime, _ := timeTransformer.ConvertStrToTime(req.Start)
	endTime, _ := timeTransformer.ConvertStrToTime(req.End)
	inboundCreated, err := h.inboundSvc.Create(&inbound.CreateInbound{
		UserID:     req.UserID,
		Start:      startTime,
		End:        endTime,
		VPNType:    entity.ConvertStringVPNTypeToEnum(req.VPNType),
		ServerType: entity.High,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "someting went wrong",
		})
	}

	return c.JSON(http.StatusAccepted, &inboundSerializer.CredateInboundSerializer{
		ID: inboundCreated.ID,
	})
}
