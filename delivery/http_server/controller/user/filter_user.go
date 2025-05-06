package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userSerializer "github.com/mohamadrezamomeni/momo/serializer/user"
)

func (u *Handler) filterUsers(c echo.Context) error {
	users, err := u.userSvc.Filter()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error please try another time",
		})
	}
	filterSerializer := &userSerializer.Filter{
		Users: make([]*userSerializer.UserSerialize, 0),
	}
	for _, user := range users {
		filterSerializer.Users = append(filterSerializer.Users, &userSerializer.UserSerialize{
			Username:     user.Username,
			IsAdmin:      user.IsAdmin,
			FirstName:    user.FirstName,
			Lastname:     user.LastName,
			IsSuperAdmin: user.IsSuperAdmin,
		})
	}
	return c.JSON(http.StatusAccepted, filterSerializer)
}
