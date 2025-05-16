package root

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) Root(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.root"

	idStr, _ := core.GetID(update)
	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(update).ErrorWrite()
	}

	titleMenu, err := telegrammessages.GetMessage("root.menu", map[string]string{})
	if err != nil {
		return nil, err
	}

	titleListVPNs, err := telegrammessages.GetMessage("inbound.list_buttom", map[string]string{})
	if err != nil {
		return nil, err
	}

	titleCreateVPNs, err := telegrammessages.GetMessage("inbound.create_buttom", map[string]string{})
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(int64(id), titleMenu)

	inlineKeyboard := [][]tgbotapi.InlineKeyboardButton{}

	if update.UserSystem != nil {
		listInboundsButtom := tgbotapi.NewInlineKeyboardButtonData(titleListVPNs, "/list_inbound")
		createInboundsButtom := tgbotapi.NewInlineKeyboardButtonData(titleCreateVPNs, "/create_inbound")
		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(listInboundsButtom, createInboundsButtom))
	} else {
		button := tgbotapi.NewInlineKeyboardButtonData("register", "/register")

		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(button))
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(inlineKeyboard...)

	msg.ReplyMarkup = markup

	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
	}, nil
}
