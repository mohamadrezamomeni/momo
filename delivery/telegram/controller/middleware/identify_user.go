package middleware

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	userServiceDto "github.com/mohamadrezamomeni/momo/dto/service/user"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type UserService interface {
	FindByTelegramID(string) (*entity.User, error)
	Create(*userServiceDto.AddUser) (*entity.User, error)
}

func IdentifyUser(userSvc UserService) core.Middleware {
	return func(next core.HandlerFunc) core.HandlerFunc {
		return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
			userTelegram, err := core.GetTelegramUser(update)
			if err != nil {
				return nil, err
			}
			userTelegramID := strconv.Itoa(int(userTelegram.ID))
			user, err := userSvc.FindByTelegramID(
				userTelegramID,
			)

			if momoErr, ok := momoError.GetMomoError(err); err != nil && (!ok || momoErr.GetErrorType() != momoError.NotFound) {
				return nil, err
			}

			if user != nil && err == nil {
				update.UserSystem = user
				return next(update)
			}

			user, err = userSvc.Create(&userServiceDto.AddUser{
				IsAdmin:          true,
				Password:         "",
				TelegramUsername: userTelegram.UserName,
				Username:         "",
				LastName:         userTelegram.LastName,
				FirstName:        userTelegram.FirstName,
				TelegramID:       userTelegramID,
			})
			if err != nil {
				return nil, err
			}

			update.UserSystem = user
			return next(update)
		}
	}
}
