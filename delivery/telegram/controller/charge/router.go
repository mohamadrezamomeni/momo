package charge

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("create_charge",
		h.CreateCharge,
		middleware.IdentifyUser(h.userSvc),
		h.SetStateCreateCharge,
	)
	telegramRouter.Register("ask_detail_charge",
		h.AskDetail,
		middleware.IdentifyUser(h.userSvc),
		h.SetStateCreateCharge,
	)
	telegramRouter.Register("answer_detail_charge",
		h.AnswerDetail,
		middleware.IdentifyUser(h.userSvc),
		h.SetStateCreateCharge,
	)
}
