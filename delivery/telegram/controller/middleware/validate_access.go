package middleware

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func ValidateAccess() core.Middleware {
	return func(next core.HandlerFunc) core.HandlerFunc {
		return func(update *core.Update) (*core.ResponseHandlerFunc, error) {
			user := update.UserSystem
			id, err := utils.ConvertToInt64(user.TelegramID)
			if err != nil {
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
			return next(update)
		}
	}
}
