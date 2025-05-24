package middleware

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
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
			id, err := utils.ConvertToInt64(idStr)
			if err != nil {
				return nil, err
			}
			user, err := userSvc.FindByTelegramID(idStr)
			if momoErr, ok := momoError.GetMomoError(err); err != nil && (!ok || momoErr.GetErrorType() != momoError.NotFound) {
				return nil, err
			}
			if user != nil && !user.IsApproved {
				message, _ := telegrammessages.GetMessage("error.forbidden_access", map[string]string{})
				msgConfig := tgbotapi.NewMessage(id, message)
				return &core.ResponseHandlerFunc{
					MessageConfig: &msgConfig,
					ReleaseState:  true,
					RedirectRoot:  false,
				}, nil
			}
			update.UserSystem = user
			return next(update)
		}
	}
}
