package vpn

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegramState "github.com/mohamadrezamomeni/momo/delivery/telegram/state"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) AskSelectingVPNType(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr, err := core.GetID(update)
	if err != nil {
		return nil, err
	}

	markup := h.getVPNTypeButtons()

	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, err
	}

	text, err := telegrammessages.GetMessage(
		"vpn.select_vpn",
		map[string]string{},
	)
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(id, text)

	msgConfig.ReplyMarkup = markup

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
	}, nil
}

func (h *Handler) getVPNTypeButtons() tgbotapi.InlineKeyboardMarkup {
	xray := tgbotapi.NewInlineKeyboardButtonData(
		entity.VPNTypeString(entity.XRAY_VPN), entity.VPNTypeString(entity.XRAY_VPN),
	)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(xray),
	)
}

func (h *Handler) SelectVPN(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.fillVPNType"
	idStr, err := core.GetID(update)
	if err != nil {
		return nil, err
	}

	state, isExist := telegramState.FindState(idStr)
	if !isExist {
		return nil, momoError.Scope(scope).ErrorWrite()
	}

	vpnTypeStr := update.CallbackQuery.Data
	state.SetData("vpn_type", vpnTypeStr)
	return nil, nil
}
