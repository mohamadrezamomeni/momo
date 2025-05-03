package scheduler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	inboundSvc InboundService
	vpnSvc     VPNService
	sch        *gocron.Scheduler
}

type InboundService interface {
	HealingUpInbounds()
	AssignDomainToInbounds()
}

type VPNService interface {
	MonitorVPNs()
}

func New(inboundSvc InboundService, vpnSvc VPNService) *Scheduler {
	return &Scheduler{
		inboundSvc: inboundSvc,
		vpnSvc:     vpnSvc,
		sch:        gocron.NewScheduler(time.UTC),
	}
}

func (s *Scheduler) Start(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Cron("*/10 * * * *").Do(s.inboundSvc.HealingUpInbounds)
	s.sch.Cron("*/5 * * * *").Do(s.inboundSvc.AssignDomainToInbounds)
	s.sch.Cron("*/2 * * * *").Do(s.vpnSvc.MonitorVPNs)

	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
