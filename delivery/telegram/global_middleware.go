package telegram

import (
	"github.com/mohamadrezamomeni/momo/delivery/telegram/controller/middleware"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
)

func (t *Telegram) getGlobalMiddlewares() []core.Middleware {
	return []core.Middleware{
		middleware.IdentifyUser(t.userSvc),
	}
}
