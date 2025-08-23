package vpnsource

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	VPNSourceServiceDto "github.com/mohamadrezamomeni/momo/dto/service/vpn_source"
	"github.com/mohamadrezamomeni/momo/entity"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) AskVPNSource(update *core.Update) (*core.ResponseHandlerFunc, error) {
	vpnSources, err := h.vpnSourceService.FilterVPNSources(&VPNSourceServiceDto.FilterVPNSourcesDto{
		Available: true,
	})
	if err != nil {
		return nil, err
	}

	if len(vpnSources) == 0 {
		return h.notFoundAnyVPNSources(update)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, vpnSource := range vpnSources {
		button := h.getVPNSourceButton(vpnSource)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(*button))
	}

	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}

	askingPackageTitle, err := telegrammessages.GetMessage(
		"vpn_source.ask_vpn_source",
		map[string]string{},
		update.UserSystem.Language,
	)
	if err != nil {
		return nil, err
	}

	msgConfig := tgbotapi.NewMessage(id, askingPackageTitle)

	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
	}, nil
}

func (h *Handler) getVPNSourceButton(vpnSource *entity.VPNSource) *tgbotapi.InlineKeyboardButton {
	button := tgbotapi.NewInlineKeyboardButtonData(vpnSource.English, vpnSource.Country)
	return &button
}
