package middleware

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	"github.com/mohamadrezamomeni/momo/entity"
)

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
}

func IdentifyUser(userSvc UserService) core.Middleware {
	return func(next core.HandlerFunc) core.HandlerFunc {
		return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
			idStr, err := core.GetID(update)
			if err != nil {
				return nil, err
			}
			user, err := userSvc.FindByTelegramID(idStr)
			if err != nil {
				return nil, err
			}
			update.UserSystem = user
			return next(update)
		}
	}
}
