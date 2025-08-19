package vpnsource

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("ask_selecting_VPNSource",
		h.AskVPNSource,
		middleware.ValidateAccess(),
	)
	telegramRouter.Register("answer_VPNSource",
		h.AnswerVPNSource,
		middleware.ValidateAccess(),
	)
}
