package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userDto "github.com/mohamadrezamomeni/momo/dto/controller/user"
	"github.com/mohamadrezamomeni/momo/dto/service/user"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
)

func (h *Handler) Create(c echo.Context) error {
	var req userDto.AddUser
	if err := c.Bind(&req); err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}

	err := h.validator.ValidateAddUserRequest(req)
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	_, err = h.userSvc.Create(&user.AddUser{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsAdmin:   req.IsAdmin,
		Password:  req.Password,
	})
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	return c.NoContent(http.StatusNoContent)
}
