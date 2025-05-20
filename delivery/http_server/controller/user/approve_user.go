package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userDto "github.com/mohamadrezamomeni/momo/dto/controller/user"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) ApproveUser(c echo.Context) error {
	var req userDto.ApproveUser
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.userSvc.ApproveUser(req.ID)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
