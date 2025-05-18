package inbound

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	httpErr "github.com/mohamadrezamomeni/momo/pkg/http_error"
	timeTransformer "github.com/mohamadrezamomeni/momo/transformer/time"
)

func (i *Handler) ExtendInbound(c echo.Context) error {
	var req inboundControllerDto.ExtendInboundDto
	if err := c.Bind(&req); err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	err := i.validation.ValidateExtendingInbound(req)
	if err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	endTime, err := timeTransformer.ConvertStrToTime(req.End)
	if err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err = i.inboundSvc.ExtendInbound(req.ID, &inboundServiceDto.ExtendInboundDto{
		ExtendedTrafficLimit: req.ExtendedTrafficLimit,
		Start:                time.Now(),
		End:                  endTime,
	})
	if err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	return c.NoContent(http.StatusOK)
}
