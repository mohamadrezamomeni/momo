package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func MenuConfigMessage(telegramID string, user *entity.User) (*tgbotapi.MessageConfig, error) {
	scope := "telegramPKG.MenuConfigMessage"
	id, err := utils.ConvertToInt64(telegramID)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(telegramID, user).ErrorWrite()
	}

	titleMenu, err := telegrammessages.GetMessage("root.menu", map[string]string{})
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(id, titleMenu)

	inlineKeyboard := [][]tgbotapi.InlineKeyboardButton{}

	inboundRow, err := getInboundRow()
	if err != nil {
		return nil, err
	}

	chargeRow, err := getChargeRow()
	if err != nil {
		return nil, err
	}

	if user != nil {
		inlineKeyboard = append(inlineKeyboard, inboundRow, chargeRow)
	} else {
		button := tgbotapi.NewInlineKeyboardButtonData("register", "/register")

		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(button))
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(inlineKeyboard...)

	msgConfig.ReplyMarkup = markup

	return &msgConfig, nil
}

func getInboundRow() ([]tgbotapi.InlineKeyboardButton, error) {
	titleListVPNs, err := telegrammessages.GetMessage("inbound.list_buttom", map[string]string{})
	if err != nil {
		return nil, err
	}

	titleCreateVPNs, err := telegrammessages.GetMessage("inbound.create_buttom", map[string]string{})
	if err != nil {
		return nil, err
	}

	listInboundsButtom := tgbotapi.NewInlineKeyboardButtonData(titleListVPNs, "/list_inbounds")
	createInboundsButtom := tgbotapi.NewInlineKeyboardButtonData(titleCreateVPNs, "/create_inbound")
	return tgbotapi.NewInlineKeyboardRow(listInboundsButtom, createInboundsButtom), nil
}

func getChargeRow() ([]tgbotapi.InlineKeyboardButton, error) {
	titleCreateCharge, err := telegrammessages.GetMessage("charge.extend_buttom", map[string]string{})
	if err != nil {
		return nil, err
	}

	createChargeButtom := tgbotapi.NewInlineKeyboardButtonData(titleCreateCharge, "/create_charge")
	return tgbotapi.NewInlineKeyboardRow(createChargeButtom), nil
}
