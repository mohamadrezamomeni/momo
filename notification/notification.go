package notification

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	eventServiceDto "github.com/mohamadrezamomeni/momo/dto/service/event"
	"github.com/mohamadrezamomeni/momo/entity"

	chargeHandler "github.com/mohamadrezamomeni/momo/notification/charge"
	"github.com/mohamadrezamomeni/momo/notification/core"
	inboundHandler "github.com/mohamadrezamomeni/momo/notification/inbound"
	userHandler "github.com/mohamadrezamomeni/momo/notification/user"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	chargeService "github.com/mohamadrezamomeni/momo/service/charge"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	userService "github.com/mohamadrezamomeni/momo/service/user"
)

type Notification struct {
	core           *core.Core
	inboundHandler *inboundHandler.Handler
	userHandler    *userHandler.Handler
	chargeHandler  *chargeHandler.Handler
	eventSvc       EventService
}

type EventService interface {
	MarkNotificationProcessed(string) error
	FilterNotifications(*eventServiceDto.FilterEvents) ([]*entity.Event, error)
}

func New(
	cfg *NotificationConfig,
	inboundSvc *inboundService.Inbound,
	userSvc *userService.User,
	chargeSvc *chargeService.Charge,
	eventSvc EventService,
) *Notification {
	scope := "notification.new"
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		momoError.Wrap(err).Scope(scope).BadRequest().Fatalf("error to initialize bot")
	}

	return &Notification{
		core:           core.New(bot),
		inboundHandler: inboundHandler.New(userSvc, inboundSvc),
		chargeHandler:  chargeHandler.New(userSvc, inboundSvc, chargeSvc),
		userHandler:    userHandler.New(userSvc),
		eventSvc:       eventSvc,
	}
}

func (n *Notification) ServeRoutes() *Notification {
	n.chargeHandler.SetRouter(n.core)
	n.inboundHandler.SetRouter(n.core)
	n.userHandler.SetRouter(n.core)
	return n
}

func (n *Notification) NotifyEvents() {
	scope := "notification.notify"
	routes := n.core.GetRoutes()

	notActive := false
	events, err := n.eventSvc.FilterNotifications(&eventServiceDto.FilterEvents{
		Names:                   routes,
		IsNotificationProcessed: &notActive,
	})
	if err != nil {
		momoError.Wrap(err).Scope(scope).DebuggingError()
		return
	}
	var wg sync.WaitGroup
	for _, event := range events {
		wg.Add(1)
		go n.notify(event, &wg)
	}
	wg.Wait()
}

func (n *Notification) notify(event *entity.Event, wg *sync.WaitGroup) {
	defer wg.Done()
	err := n.core.Notify(event.Name, event.Data)
	if err == nil {
		n.eventSvc.MarkNotificationProcessed(event.ID)
	}
}
