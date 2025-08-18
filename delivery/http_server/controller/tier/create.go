package tier

import (
	"net/http"

	"github.com/labstack/echo/v4"
	tierControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/tier"
	tierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/tier"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) Create(c echo.Context) error {
	var req tierControllerDto.CreateTier
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	if err := h.tierValiator.ValidateCreatingTier(req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	isDefault := false
	if req.IsDefault != nil {
		isDefault = *req.IsDefault
	}
	_, err := h.tierSvc.Create(&tierServiceDto.CreateTier{
		IsDefault: isDefault,
		Name:      req.Name,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
