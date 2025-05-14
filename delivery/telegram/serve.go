package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	authService "github.com/mohamadrezamomeni/momo/service/auth"
	userService "github.com/mohamadrezamomeni/momo/service/user"

	authHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/auth"
	rootHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/root"
)

type Telegram struct {
	config      *TelegramConfig
	userSvc     *userService.User
	bot         *tgbotapi.BotAPI
	core        *core.Router
	authHandler *authHandler.Handler
	rootHandler *rootHandler.Handler
}

func New(cfg *TelegramConfig, userSvc *userService.User, authSvc *authService.Auth) *Telegram {
	scope := "telegram.New"
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		momoError.Wrap(err).Scope(scope).BadRequest().Fatalf("error to initialize bot")
	}

	rootHandler := rootHandler.New(userSvc)
	return &Telegram{
		rootHandler: rootHandler,
		bot:         bot,
		core:        core.New("menu"),
		config:      cfg,
		userSvc:     userSvc,
		authHandler: authHandler.New(authSvc, userSvc),
	}
}

func (t *Telegram) Serve() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	t.authHandler.SetRouter(t.core)
	t.rootHandler.SetRouter(t.core)

	for update := range updates {
		res := t.core.Route(&update)
		if res != nil {
			t.send(res, &update)
		}
	}
}

func (t *Telegram) send(res *core.ResponseHandlerFunc, update *tgbotapi.Update) {
	t.bot.Send(res.Result)

	if res.RedirectRoot {
		r, _ := t.rootHandler.Root(update)
		t.bot.Send(r.Result)
	}
}
