package vpnpackage

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("ask_selecting_VPNPackage", h.AskSelectingVPNPackage, middleware.IdentifyUser(h.userSvc))
	telegramRouter.Register("answer_selecting_VPNPackage", h.SelectVPNPackage, middleware.IdentifyUser(h.userSvc))
}
