package scheduler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	healingUpInboundSvc HealingUpInboundService
	inboundTrafficSvc   InboundTrafficService
	hostInboundSvc      HostInboundService
	vpnSvc              VPNService
	hostSvc             HostService
	sch                 *gocron.Scheduler
	inboundChargeSvc    InboundChargeService
}

type HealingUpInboundService interface {
	HealingUpExpiredInbounds()
	HealingUpOverQuotedInbounds()
	HealingUpBlockedInbounds()
	HealingUpChargedInbounds()
	CheckInboundsActivation()
}

type InboundTrafficService interface {
	UpdateTraffics()
}

type InboundChargeService interface {
	ChargeInbounds()
}

type HostInboundService interface {
	AssignDomainToInbounds()
}

type VPNService interface {
	MonitorVPNs()
}

type HostService interface {
	MonitorHosts()
}

func New(
	healingUpInboundSvc HealingUpInboundService,
	inboundTrafficSvc InboundTrafficService,
	hostInboundSvc HostInboundService,
	vpnSvc VPNService,
	hostSvc HostService,
	inboundChargeSvc InboundChargeService,
) *Scheduler {
	return &Scheduler{
		inboundChargeSvc:    inboundChargeSvc,
		healingUpInboundSvc: healingUpInboundSvc,
		inboundTrafficSvc:   inboundTrafficSvc,
		hostInboundSvc:      hostInboundSvc,
		vpnSvc:              vpnSvc,
		hostSvc:             hostSvc,
		sch:                 gocron.NewScheduler(time.UTC),
	}
}

func (s *Scheduler) Start(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	s.sch.Cron("*/5 * * * *").Do(s.healingUpInboundSvc.CheckInboundsActivation)
	s.sch.Cron("*/10 * * * *").Do(s.healingUpInboundSvc.HealingUpExpiredInbounds)
	s.sch.Cron("*/10 * * * *").Do(s.healingUpInboundSvc.HealingUpOverQuotedInbounds)
	s.sch.Cron("*/10 * * * *").Do(s.healingUpInboundSvc.HealingUpBlockedInbounds)
	s.sch.Cron("*/10 * * * *").Do(s.healingUpInboundSvc.HealingUpChargedInbounds)
	s.sch.Cron("*/5 * * * *").Do(s.hostInboundSvc.AssignDomainToInbounds)
	s.sch.Cron("*/1 * * * *").Do(s.inboundTrafficSvc.UpdateTraffics)
	s.sch.Cron("*/2 * * * *").Do(s.vpnSvc.MonitorVPNs)
	s.sch.Cron("*/2 * * * *").Do(s.hostSvc.MonitorHosts)
	s.sch.Cron("*/5 * * * *").Do(s.inboundChargeSvc.ChargeInbounds)
	s.sch.StartAsync()

	<-done
	s.sch.Stop()
}
