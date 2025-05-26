package user

import "github.com/mohamadrezamomeni/momo/notification/core"

func (h *Handler) SetRouter(notificationRouter *core.Core) {
	notificationRouter.Register("approve_user", h.ApproveUserByAddmin)
}
