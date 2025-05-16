package middleware

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	userService "github.com/mohamadrezamomeni/momo/service/user"
)

func IdentifyUser(userSvc *userService.User) core.Middleware {
	return func(next core.HandlerFunc) core.HandlerFunc {
		return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
			idStr, err := core.GetID(update)
			if err != nil {
				return nil, err
			}
			user, err := userSvc.FindByUsername(idStr)
			if err != nil {
				return nil, err
			}
			update.UserSystem = user
			return next(update)
		}
	}
}
