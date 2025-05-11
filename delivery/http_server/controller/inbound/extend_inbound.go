package inbound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	inboundControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/inbound"
	httpErr "github.com/mohamadrezamomeni/momo/pkg/http_error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
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

	err = i.inboundSvc.ExtendInbound(req.ID, utils.GetDateTime(req.End))
	if err != nil {
		msg, code := httpErr.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	return c.NoContent(http.StatusOK)
}
