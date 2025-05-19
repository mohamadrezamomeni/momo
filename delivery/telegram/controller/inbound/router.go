package inbound

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("list_inbounds", h.ListInbounds, middleware.IdentifyUser(h.userSvc))
	telegramRouter.Register("create_inbound",
		h.CreateInbound,
		middleware.IdentifyUser(h.userSvc),
		h.SetState,
		h.AskVPNType,
		h.FillVPNType,
		h.AskPackagesCreatingInbound,
		h.AnswerPackageCreatingInbound,
	)
	telegramRouter.Register("extend_inbound", h.ExtendInbound,
		middleware.IdentifyUser(h.userSvc),
		h.SetExtendingInboundState,
		h.SelectInboundIDInExtendingInbound,
		h.ChooseInbound,
		h.AskPackagesExtending,
		h.AnswerPackageExtendingInbound,
	)
}
