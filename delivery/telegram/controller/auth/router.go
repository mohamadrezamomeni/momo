package auth

import "github.com/mohamadrezamomeni/momo/delivery/telegram/core"

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("register",
		h.Register,
		h.SetState,
	)
	telegramRouter.Register("askUsername",
		h.AskUsername,
	)
	telegramRouter.Register("answerUsername",
		h.AnswerUsername,
	)
}
