package root

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) Root(update *tgbotapi.Update) (*core.ResponseHandlerFunc, error) {
	button := tgbotapi.NewInlineKeyboardButtonData("register", "register")

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please press the button:")
	msg.ReplyMarkup = markup

	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
	}, nil
}
