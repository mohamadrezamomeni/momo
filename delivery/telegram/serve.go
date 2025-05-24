package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/momo/delivery/telegram/core"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	telegrammessages "github.com/mohamadrezamomeni/momo/pkg/telegram_messages"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
	authService "github.com/mohamadrezamomeni/momo/service/auth"
	chargeService "github.com/mohamadrezamomeni/momo/service/charge"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	vpnPackageService "github.com/mohamadrezamomeni/momo/service/vpn_package"

	authHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/auth"
	chargeHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/charge"
	inboundHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/inbound"
	rootHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/root"
	vpnHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/vpn"
	vpnPackageHandler "github.com/mohamadrezamomeni/momo/delivery/telegram/controller/vpn_package"

	inboundValidator "github.com/mohamadrezamomeni/momo/validator/inbound"
)

type Telegram struct {
	config            *TelegramConfig
	userSvc           *userService.User
	bot               *tgbotapi.BotAPI
	core              *core.Router
	authHandler       *authHandler.Handler
	rootHandler       *rootHandler.Handler
	inboundHandler    *inboundHandler.Handler
	vpnHandler        *vpnHandler.Handler
	vpnPackageHandler *vpnPackageHandler.Handler
	chargeHandler     *chargeHandler.Handler
}

func New(
	cfg *TelegramConfig,
	userSvc *userService.User,
	authSvc *authService.Auth,
	inboundSvc *inboundService.Inbound,
	vpnPackageSvc *vpnPackageService.VPNPackage,
	chargeSvc *chargeService.Charge,
) *Telegram {
	scope := "telegram.New"
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		momoError.Wrap(err).Scope(scope).BadRequest().Fatalf("error to initialize bot")
	}

	return &Telegram{
		bot:               bot,
		core:              core.New("menu"),
		config:            cfg,
		userSvc:           userSvc,
		rootHandler:       rootHandler.New(userSvc),
		authHandler:       authHandler.New(authSvc, userSvc),
		vpnHandler:        vpnHandler.New(userSvc),
		vpnPackageHandler: vpnPackageHandler.New(vpnPackageSvc, userSvc),
		chargeHandler:     chargeHandler.New(chargeSvc, userSvc),
		inboundHandler: inboundHandler.New(
			userSvc,
			inboundSvc,
			vpnPackageSvc,
			inboundValidator.New(userSvc, inboundSvc),
		),
	}
}

func (t *Telegram) Serve() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	t.authHandler.SetRouter(t.core)
	t.rootHandler.SetRouter(t.core)
	t.inboundHandler.SetRouter(t.core)
	t.vpnHandler.SetRouter(t.core)
	t.vpnPackageHandler.SetRouter(t.core)
	t.chargeHandler.SetRouter(t.core)

	for update := range updates {
		customUpdate := &core.Update{
			Update: &update,
		}

		res, err := t.core.Route(customUpdate)
		if err != nil {
			t.sendError(customUpdate)
		}

		if res != nil {
			t.send(res, customUpdate)
		}
	}
}

func (t *Telegram) send(res *core.ResponseHandlerFunc, update *core.Update) {
	t.bot.Send(res.MessageConfig)

	if res.RedirectRoot {
		r, _ := t.core.RootHandler(update)
		t.bot.Send(r.MessageConfig)
	}
}

func (t *Telegram) sendError(update *core.Update) {
	errMessage, _ := telegrammessages.GetMessage("error.internal_error", map[string]string{})
	idStr, _ := core.GetID(update)
	id, _ := utils.ConvertToInt64(idStr)
	msg := tgbotapi.NewMessage(id, errMessage)
	t.bot.Send(msg)
}
