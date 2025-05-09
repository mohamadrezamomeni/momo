package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/auth"
	authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	authSerializer "github.com/mohamadrezamomeni/momo/serializer/auth"
)

func (h *Handler) Login(c echo.Context) error {
	var req authControllerDto.Login
	if err := c.Bind(&req); err != nil {
		message, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": message,
		})
	}
	err := h.validation.LoginValidator(req)
	if err != nil {
		message, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": message,
		})
	}
	accessToken, _, err := h.authSvc.Login(&authServiceDto.LoginDto{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		message, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": message,
		})
	}

	return c.JSON(http.StatusAccepted, &authSerializer.Login{
		AccessToken: accessToken,
	})
}
