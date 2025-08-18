package usertier

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userTierControllereDto "github.com/mohamadrezamomeni/momo/dto/controller/user_tier"
	userTierServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user_tier"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) Create(c echo.Context) error {
	var req userTierControllereDto.Create
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	if err := h.userTierValidator.ValidateCreatingUserTier(req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.userTierSvc.Create(&userTierServiceDto.Create{
		UserID: req.UserID,
		Tier:   req.Tier,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
