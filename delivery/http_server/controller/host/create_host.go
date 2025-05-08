package host

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hostControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/host"
	hostServiceDto "github.com/mohamadrezamomeni/momo/dto/service/host"
)

func (h *Handler) CreateHost(c echo.Context) error {
	var req hostControllerDto.CreateHostDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "input is wrong",
		})
	}

	err := h.hostValidator.CreateHostValidation(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "input is wrong",
		})
	}

	_, err = h.hostSvc.Create(&hostServiceDto.CreateHostDto{
		Domain: req.Domain,
		Port:   req.Port,
	})
	return c.NoContent(http.StatusNoContent)
}
