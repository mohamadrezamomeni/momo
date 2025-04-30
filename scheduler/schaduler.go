package scheduler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	inboundSvc InboundService
	sch        *gocron.Scheduler
}

type InboundService interface {
	HealingUpInbounds()
}

func New(inboundSvc InboundService) *Scheduler {
	return &Scheduler{
		inboundSvc: inboundSvc,
		sch:        gocron.NewScheduler(time.UTC),
	}
}

func (s *Scheduler) Start(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Cron("*/10 * * * *").Do(s.inboundSvc.HealingUpInbounds)

	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
