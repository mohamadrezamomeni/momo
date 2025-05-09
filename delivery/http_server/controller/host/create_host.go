package host

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hostControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/host"
	hostServiceDto "github.com/mohamadrezamomeni/momo/dto/service/host"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) CreateHost(c echo.Context) error {
	var req hostControllerDto.CreateHostDto
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.hostValidator.CreateHostValidation(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	_, err = h.hostSvc.Create(&hostServiceDto.CreateHostDto{
		Domain: req.Domain,
		Port:   req.Port,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
