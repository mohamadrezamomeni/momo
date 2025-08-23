package middleware

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func ValidateAccess() core.Middleware {
	scope := "telegram.middleware.ValidateAccess"
	return func(next core.HandlerFunc) core.HandlerFunc {
		return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
			user := update.UserSystem
			if user == nil {
				return nil, momoError.Scope(scope).Errorf("the user is nill")
			}
			id, err := utils.ConvertToInt64(user.TelegramID)
			if err != nil {
				return nil, err
			}
			if user != nil && !user.IsApproved {
				message, _ := telegrammessages.GetMessage(
					"error.forbidden_access",
					map[string]string{},
					update.UserSystem.Language,
				)
				msgConfig := tgbotapi.NewMessage(id, message)
				return &core.ResponseHandlerFunc{
					MessageConfig: &msgConfig,
					ReleaseState:  true,
					RedirectRoot:  false,
				}, nil
			}
			return next(update)
		}
	}
}
