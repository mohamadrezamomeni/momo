package root

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func (h *Handler) Start(update *core.Update) (*core.ResponseHandlerFunc, error) {
	scope := "telegram.controller.start"

	idStr, _ := core.GetID(update)
	id, err := utils.ConvertToInt64(idStr)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(update).ErrorWrite()
	}

	msgConfig := tgbotapi.NewMessage(int64(id), "hello welcome to our home")

	return &core.ResponseHandlerFunc{
		MessageConfig: &msgConfig,
		ReleaseState:  true,
		RedirectRoot:  true,
	}, nil
}
