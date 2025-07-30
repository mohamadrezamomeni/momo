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

	clienConfigRow, err := getClientConfigRow()
	if err != nil {
		return nil, err
	}

	if user != nil {
		inlineKeyboard = append(inlineKeyboard, chargeRow, inboundRow, clienConfigRow)
	} else {
		button := tgbotapi.NewInlineKeyboardButtonData("register", "/register")

		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(button))
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(inlineKeyboard...)

	msgConfig.ReplyMarkup = markup

	return &msgConfig, nil
}

func getClientConfigRow() ([]tgbotapi.InlineKeyboardButton, error) {
	titleGenerateClientConfig, err := telegrammessages.GetMessage("inbound.client_config_button", map[string]string{})
	if err != nil {
		return nil, err
	}
	generateClientConfig := tgbotapi.NewInlineKeyboardButtonData(titleGenerateClientConfig, "/generate_client_config")
	return tgbotapi.NewInlineKeyboardRow(generateClientConfig), nil
}

func getInboundRow() ([]tgbotapi.InlineKeyboardButton, error) {
	titleListVPNs, err := telegrammessages.GetMessage("inbound.list_button", map[string]string{})
	if err != nil {
		return nil, err
	}

	listInboundsButton := tgbotapi.NewInlineKeyboardButtonData(titleListVPNs, "/list_inbounds")
	return tgbotapi.NewInlineKeyboardRow(listInboundsButton), nil
}

func getChargeRow() ([]tgbotapi.InlineKeyboardButton, error) {
	titleChargeInbound, err := telegrammessages.GetMessage("charge.extend_button", map[string]string{})
	if err != nil {
		return nil, err
	}
	titleCreateVPN, err := telegrammessages.GetMessage("inbound.create_button", map[string]string{})
	if err != nil {
		return nil, err
	}
	createInboundsButton := tgbotapi.NewInlineKeyboardButtonData(titleCreateVPN, "/create_charge")

	createChargeButton := tgbotapi.NewInlineKeyboardButtonData(titleChargeInbound, "/charge_inbound")
	return tgbotapi.NewInlineKeyboardRow(createChargeButton, createInboundsButton), nil
}
