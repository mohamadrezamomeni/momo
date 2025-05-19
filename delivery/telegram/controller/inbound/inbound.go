package inbound

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	"github.com/mohamadrezamomeni/momo/entity"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) sendNotFoundInbounds(user *entity.User) (*core.ResponseHandlerFunc, error) {
	id, err := utils.ConvertToInt64(user.TelegramID)
	if err != nil {
		return nil, err
	}
	inboundNotFoundText, _ := telegrammessages.GetMessage("inbound.not_found", map[string]string{})

	msg := tgbotapi.NewMessage(id, inboundNotFoundText)
	return &core.ResponseHandlerFunc{
		Result:       msg,
		ReleaseState: true,
		RedirectRoot: true,
	}, nil
}
