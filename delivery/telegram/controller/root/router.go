package root

import "github.com/mohamadrezamomeni/momo/delivery/telegram/core"

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("start", h.Start)
	telegramRouter.Register("menu", h.Root)
}
