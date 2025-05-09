package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) Block(c echo.Context) error {
	req := inboundControllerDto.IdentifyInbounbdDto{}
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.inboundSvc.Block(req.ID)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) UnBlock(c echo.Context) error {
	req := inboundControllerDto.IdentifyInbounbdDto{}
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.inboundSvc.UnBlock(req.ID)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusOK)
}
