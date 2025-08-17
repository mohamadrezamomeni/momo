package tier

import (
	"net/http"

	"github.com/labstack/echo/v4"
	tierControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/tier"
	tierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/tier"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) Update(c echo.Context) error {
	var req tierControllerDto.Update
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.tierSvc.Update(req.Name, &tierServiceDto.Update{
		IsDefault: req.IsDefault,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
