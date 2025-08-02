package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	momoErrorHttp "github.com/mohamadrezamomeni/momo/pkg/http_error"
	userSerializer "github.com/mohamadrezamomeni/momo/serializer/user"
)

func (u *Handler) filterUsers(c echo.Context) error {
	users, err := u.userSvc.Filter()
	if err != nil {
		msg, code := momoErrorHttp.Error(err)
		return c.JSON(code, map[string]string{
			"message": msg,
		})
	}
	filterSerializer := &userSerializer.Filter{
		Users: make([]*userSerializer.UserSerialize, 0),
	}
	for _, user := range users {
		filterSerializer.Users = append(filterSerializer.Users, &userSerializer.UserSerialize{
			ID:           user.ID,
			Username:     user.Username,
			IsAdmin:      user.IsAdmin,
			FirstName:    user.FirstName,
			Lastname:     user.LastName,
			IsSuperAdmin: user.IsSuperAdmin,
			IsApproved:   user.IsApproved,
		})
	}
	return c.JSON(http.StatusAccepted, filterSerializer)
}
