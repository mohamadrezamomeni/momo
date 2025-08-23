package vpnpackage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) notFoundAnyVPNPackages(update *core.Update) (*core.ResponseHandlerFunc, error) {
	text, err := telegrammessages.GetMessage(
		"vpn_package.not_found_package",
		map[string]string{},
		update.UserSystem.Language,
	)
	if err != nil {
		return nil, err
	}
	id, err := utils.ConvertToInt64(update.UserSystem.TelegramID)
	if err != nil {
		return nil, err
	}
	msgConfig := tgbotapi.NewMessage(id, text)
	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		MenuTab:       true,
		ReleaseState:  true,
		RedirectRoot:  true,
	}, nil
}
