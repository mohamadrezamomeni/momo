package root

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"

	telegramPKG "github.com/mohamadrezamomeni/momo/pkg/telegram"
)

func (h *Handler) Root(update *core.Update) (*core.ResponseHandlerFunc, error) {
	idStr, _ := core.GetID(update)

	msgConfig, err := telegramPKG.MenuConfigMessage(idStr, update.UserSystem)
	if err != nil {
		return nil, err
	}

	return &core.ResponseHandlerFunc{
		MessageConfig: msgConfig,
		ReleaseState:  true,
	}, nil
}
