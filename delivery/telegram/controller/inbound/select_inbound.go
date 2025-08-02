package inbound

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) AskSelectingInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	inbounds, err := h.inboundSvc.Filter(&inboundServiceDto.FilterInbounds{
		UserID: update.UserSystem.ID,
	})

	if len(inbounds) == 0 {
		return h.sendNotFoundInbounds(update.UserSystem)
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}
	askInboundID, err := telegrammessages.GetMessage("inbound.select_vpn.ask_id", map[string]string{})
	if err != nil {
		return nil, err
	}
	msgConfig := tgbotapi.NewMessage(id, askInboundID)
	var sb strings.Builder

	sb.WriteString(askInboundID)

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, inbound := range inbounds {
		button, err := h.makeExtendingInboundButton(inbound)
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

func (h *Handler) makeExtendingInboundButton(inbound *entity.Inbound) (*tgbotapi.InlineKeyboardButton, error) {
	itemTtitle, err := telegrammessages.GetMessage("inbound.select_vpn.item", map[string]string{
		"VPNType": entity.VPNTypeString(inbound.VPNType),
		"ID":      strconv.Itoa(inbound.ID),
	})
	if err != nil {
		return nil, err
	}
	button := tgbotapi.NewInlineKeyboardButtonData(itemTtitle, strconv.Itoa(inbound.ID))
	return &button, nil
}

func (h *Handler) AnswerSelectingInbound(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.inbound.chooseInbound"
	inboundID := update.CallbackQuery.Data

	inbound, err := h.inboundSvc.FindInboundByID(inboundID)
	if err != nil {
		return nil, err
	}
	err = h.inboundValidator.ValidateExtendingInboundByUser(inbound, update.UserSystem)
	if err != nil {
		return nil, err
	}
	state, isExist := telegramState.FindState(update.UserSystem.TelegramID)

	if !isExist {
		return nil, momoError.Scope(scope).ErrorWrite()
	}

	state.SetData("inbound_id", strconv.Itoa(inbound.ID))
	return nil, nil
}
