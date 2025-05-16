package root

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("start", h.Start, middleware.IdentifyUser(h.userSvc))
	telegramRouter.Register("menu", h.Root, middleware.IdentifyUser(h.userSvc))
}
