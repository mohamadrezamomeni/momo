package charge

import "github.com/mohamadrezamomeni/momo/notification/core"

func (h *Handler) SetRouter(notificationRouter *core.Core) {
	notificationRouter.Register("charge_approve", h.ApproveCharge)
}
