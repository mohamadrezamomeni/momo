package scheduler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	inboundSvc InboundService
	vpnSvc     VPNService
	hostSvc    HostService
	sch        *gocron.Scheduler
}

type InboundService interface {
	HealingUpInbounds()
	AssignDomainToInbounds()
}

type VPNService interface {
	MonitorVPNs()
}

type HostService interface {
	MonitorHosts()
}

func New(inboundSvc InboundService, vpnSvc VPNService, hostSvc HostService) *Scheduler {
	return &Scheduler{
		inboundSvc: inboundSvc,
		vpnSvc:     vpnSvc,
		hostSvc:    hostSvc,
		sch:        gocron.NewScheduler(time.UTC),
	}
}

func (s *Scheduler) Start(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Cron("*/10 * * * *").Do(s.inboundSvc.HealingUpInbounds)
	s.sch.Cron("*/5 * * * *").Do(s.inboundSvc.AssignDomainToInbounds)
	s.sch.Cron("*/2 * * * *").Do(s.vpnSvc.MonitorVPNs)
	s.sch.Cron("*/2 * * * *").Do(s.hostSvc.MonitorHosts)
	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
