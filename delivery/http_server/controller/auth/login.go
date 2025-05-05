package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authControllerDto "github.com/mohamadrezamomeni/momo/dto/controller/auth"
	authServiceDto "github.com/mohamadrezamomeni/momo/dto/service/auth"
	authSerializer "github.com/mohamadrezamomeni/momo/serializer/auth"
)

func (h *Handler) Login(c echo.Context) error {
	var req authControllerDto.Login
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "fields that were given was wrong",
		})
	}
	err := h.validation.LoginValidator(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "fields that were given was wrong",
		})
	}
	accessToken, _, err := h.authSvc.Login(&authServiceDto.LoginDto{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "someting went wrong",
		})
	}

	return c.JSON(http.StatusAccepted, &authSerializer.Login{
		AccessToken: accessToken,
	})
}
