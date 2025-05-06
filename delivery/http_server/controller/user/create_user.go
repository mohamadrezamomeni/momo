package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userDto "github.com/mohamadrezamomeni/momo/dto/controller/user"
	"github.com/mohamadrezamomeni/momo/dto/service/user"
)

func (h *Handler) Create(c echo.Context) error {
	var req userDto.AddUser
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err := h.validator.ValidateAddUserRequest(req)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	_, err = h.userSvc.Create(&user.AddUser{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsAdmin:   req.IsAdmin,
		Password:  req.Password,
	})
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusNoContent)
}
