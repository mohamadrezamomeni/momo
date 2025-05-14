package auth

import "github.com/mohamadrezamomeni/momo/delivery/telegram/core"

func (h *Handler) SetRouter(telegramRouter *core.Router) {
	telegramRouter.Register("register", h.Register, h.SetUserRegisteration, h.SetUsername, h.SetFirstname, h.SetLastname)
}
