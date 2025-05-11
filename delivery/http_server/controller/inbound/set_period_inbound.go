package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	transformer "github.com/mohamadrezamomeni/momo/transformer/time"
)

func (h *Handler) SetPeriodInbound(c echo.Context) error {
	var req inboundControllerDto.SetPeriodDto
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.validation.ValidateSettingPeriodTime(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	startTime, _ := transformer.ConvertStrToTime(req.Start)
	endTime, _ := transformer.ConvertStrToTime(req.End)
	err = h.inboundSvc.SetPeriodTime(req.ID, &inboundServiceDto.SetPeriodDto{
		Start: startTime,
		End:   endTime,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusOK)
}
