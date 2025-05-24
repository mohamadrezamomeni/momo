package inbound

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("list_inbounds",
		h.ListInbounds,
		middleware.IdentifyUser(h.userSvc),
	)
	telegramRouter.Register("create_inbound",
		h.CreateInbound,
		middleware.IdentifyUser(h.userSvc),
		h.SetStateCreateInbound,
	)
	telegramRouter.Register("ask_selecting_inbound",
		h.AskSelectingInbound,
		middleware.IdentifyUser(h.userSvc),
	)
	telegramRouter.Register("answer_selecting_inbound",
		h.AnswerSelectingInbound,
		middleware.IdentifyUser(h.userSvc),
	)
}
