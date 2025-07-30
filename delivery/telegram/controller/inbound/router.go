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
	telegramRouter.Register("ask_selecting_inbound",
		h.AskSelectingInbound,
		middleware.IdentifyUser(h.userSvc),
	)
	telegramRouter.Register("answer_selecting_inbound",
		h.AnswerSelectingInbound,
		middleware.IdentifyUser(h.userSvc),
	)

	telegramRouter.Register("generate_client_config",
		h.GetClientConfig,
		middleware.IdentifyUser(h.userSvc),
		h.SetStateGettingClientConfig,
	)

	telegramRouter.Register("render_client_config_buttons",
		h.RenderClientConfigButtons,
		middleware.IdentifyUser(h.userSvc),
	)
}
