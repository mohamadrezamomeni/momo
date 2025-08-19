package inbound

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("list_inbounds",
		h.ListInbounds,
		middleware.ValidateAccess(),
	)
	telegramRouter.Register("ask_selecting_inbound",
		h.AskSelectingInbound,
		middleware.ValidateAccess(),
	)
	telegramRouter.Register("answer_selecting_inbound",
		h.AnswerSelectingInbound,
		middleware.ValidateAccess(),
	)

	telegramRouter.Register("generate_client_config",
		h.GetClientConfig,
		middleware.ValidateAccess(),
		h.SetStateGettingClientConfig,
	)

	telegramRouter.Register("render_client_config_buttons",
		h.RenderClientConfigButtons,
		middleware.ValidateAccess(),
	)
}
