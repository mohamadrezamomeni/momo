package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
)

func (h *Handler) Block(c echo.Context) error {
	req := inboundControllerDto.IdentifyInbounbdDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
		})
	}

	err := h.inboundSvc.Block(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "something went wrong",
		})
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) UnBlock(c echo.Context) error {
	req := inboundControllerDto.IdentifyInbounbdDto{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "the input was wrong",
		})
	}

	err := h.inboundSvc.UnBlock(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "something went wrong",
		})
	}
	return c.NoContent(http.StatusOK)
}
