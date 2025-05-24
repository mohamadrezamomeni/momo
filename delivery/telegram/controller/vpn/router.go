package vpn

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("ask_selecting_VPN", h.AskSelectingVPNType, middleware.IdentifyUser(h.userSvc))
	telegramRouter.Register("answer_selecting_VPN", h.SelectVPN, middleware.IdentifyUser(h.userSvc))
}
