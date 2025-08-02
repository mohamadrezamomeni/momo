package charge

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("create_charge",
		h.CreateInbound,
		middleware.IdentifyUser(h.userSvc),
		middleware.ValidateAccess(),
		h.SetStateCreateCharge,
	)
	telegramRouter.Register("ask_detail_charge",
		h.AskDetail,
		middleware.IdentifyUser(h.userSvc),
		middleware.ValidateAccess(),
	)
	telegramRouter.Register("answer_detail_charge",
		h.AnswerDetail,
		middleware.IdentifyUser(h.userSvc),
		middleware.ValidateAccess(),
	)
	telegramRouter.Register("charge_inbound",
		h.ChargeInbound,
		middleware.IdentifyUser(h.userSvc),
		middleware.ValidateAccess(),
		h.SetStateChargeInbound,
	)
}
