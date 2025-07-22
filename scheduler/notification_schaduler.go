package scheduler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type NotificationScheduler struct {
	notificationSvc NotificationService
	sch             *gocron.Scheduler
}

type NotificationService interface {
	NotifyEvents()
}

func NewNotificationSchaduler(
	notificationSvc NotificationService,
) *NotificationScheduler {
	return &NotificationScheduler{
		sch:             gocron.NewScheduler(time.UTC),
		notificationSvc: notificationSvc,
	}
}

func (ns *NotificationScheduler) Start(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ns.sch.Cron("*/1 * * * *").Do(ns.notificationSvc.NotifyEvents)
	ns.sch.StartAsync()

	<-done
	ns.sch.Stop()
}
