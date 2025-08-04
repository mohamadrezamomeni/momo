package inbound

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) RenderClientConfigButtons(update *core.Update) (*core.ResponseHandlerFunc, error) {
	inbounds, err := h.inboundSvc.Filter(&inboundServiceDto.FilterInbounds{
		UserID: update.UserSystem.ID,
	})
	if err != nil {
		return nil, err
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}

	title, err := telegrammessages.GetMessage("inbound.client_config.title_list_inbound", map[string]string{})
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(id, title)

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, inbound := range inbounds {
		button, err := h.renderClientConfigButton(inbound)
		if err != nil {
			return nil, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(*button))
	}

	msgConfig.ParseMode = "HTML"
	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
	}, nil
}

func (h *Handler) renderClientConfigButton(inbound *entity.Inbound) (*tgbotapi.InlineKeyboardButton, error) {
	blockedTitle, err := telegrammessages.GetMessage("inbound.client_config.client_config_item_block", map[string]string{
		"VPNType": entity.VPNTypeString(inbound.VPNType),
		"ID":      inbound.ID,
	})
	if err != nil {
		return nil, err
	}

	okConfigItem, err := telegrammessages.GetMessage("inbound.client_config.client_config_item_ok", map[string]string{
		"VPNType": entity.VPNTypeString(inbound.VPNType),
		"ID":      inbound.ID,
	})
	if err != nil {
		return nil, err
	}

	pendingItem, err := telegrammessages.GetMessage("inbound.client_config.client_config_item_pending", map[string]string{
		"VPNType": entity.VPNTypeString(inbound.VPNType),
		"ID":      inbound.ID,
	})
	if err != nil {
		return nil, err
	}
	var button tgbotapi.InlineKeyboardButton
	if inbound.IsBlock {
		button = tgbotapi.NewInlineKeyboardButtonData(blockedTitle, inbound.ID)
	} else if !inbound.IsAssigned {
		button = tgbotapi.NewInlineKeyboardButtonData(pendingItem, inbound.ID)
	} else {
		button = tgbotapi.NewInlineKeyboardButtonData(okConfigItem, inbound.ID)
	}
	return &button, nil
}
