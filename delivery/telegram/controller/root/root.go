package root

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) Root(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.root"

	idStr, _ := core.GetID(update)
	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(update).ErrorWrite()
	}

	msg := tgbotapi.NewMessage(int64(id), "Please press the button:")

	button := tgbotapi.NewInlineKeyboardButtonData("register", "/register")
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)

	msg.ReplyMarkup = markup

	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
	}, nil
}
