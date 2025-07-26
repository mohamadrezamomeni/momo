package inbound

import (
	hostServiceDto "github.com/mohamadrezamomeni/momo/dto/service/host"
	"github.com/mohamadrezamomeni/momo/entity"
)

type HostInbound struct {
	inboundRepo   InboundHostRepo
	hostService   HostService
	vpnManagerSvc VPNManagerService
}

type VPNManagerService interface {
	GetAvailableVPNSourceDomains(vpnsSources []string) (map[string][]string, error)
}

type InboundHostRepo interface {
	FindInboundIsNotAssigned() ([]*entity.Inbound, error)
	GetListOfPortsByDomain() ([]struct {
		Domain string
		Ports  []string
	}, error)
	UpdateDomainPort(int, string, string) error
}

type HostService interface {
	ResolveHostPortPair(map[string][]string, map[string]uint32) (
		map[string][]string,
		error,
	)
}

type Address struct {
	Domain string
	Port   string
}

func NewHostInbound(
	inboundRepo InboundHostRepo,
	hostService HostService,
	VPNManagerSvc VPNManagerService,
) *HostInbound {
	return &HostInbound{
		inboundRepo:   inboundRepo,
		hostService:   hostService,
		vpnManagerSvc: VPNManagerSvc,
	}
}

func (i *HostInbound) AssignDomainToInbounds() {
	inbounds, err := i.inboundRepo.FindInboundIsNotAssigned()
	if err != nil {
		return
	}
	portUserSummery, err := i.summeryDomainPorts()
	if err != nil {
		return
	}
	vpnSources := i.getVPNSourcesFromInbounds(inbounds)

	VPNSourceDomains, err := i.vpnManagerSvc.GetAvailableVPNSourceDomains(vpnSources)
	if err != nil {
		return
	}

	hostPortPairsMap, err := i.hostService.ResolveHostPortPair(
		portUserSummery,
		i.countRequiredPortEachHost(inbounds, VPNSourceDomains),
	)
	if err != nil {
		return
	}

	VPNSourceAddresses := i.getVPNSourceInboundDestination(
		VPNSourceDomains,
		hostPortPairsMap,
	)
	seen := map[string]uint32{}
	for _, inbound := range inbounds {
		count := seen[inbound.Country]

		if len(VPNSourceAddresses[inbound.Country]) > int(count) {
			inboundDestination := VPNSourceAddresses[inbound.Country][count]
			i.inboundRepo.UpdateDomainPort(inbound.ID,
				inboundDestination.Domain,
				inboundDestination.Port,
			)
			seen[inbound.Country] += 1
		}
	}
}

func (i *HostInbound) summeryDomainPorts() (map[string][]string, error) {
	summery, err := i.inboundRepo.GetListOfPortsByDomain()
	if err != nil {
		return nil, err
	}
	res := map[string][]string{}
	for _, item := range summery {
		res[item.Domain] = item.Ports
	}
	return res, nil
}

func (i *HostInbound) getVPNSourcesFromInbounds(inbounds []*entity.Inbound) []string {
	vpnSources := []string{}
	seen := map[string]struct{}{}
	for _, inbound := range inbounds {
		if _, ok := seen[inbound.Country]; !ok {
			seen[inbound.Country] = struct{}{}
			vpnSources = append(vpnSources, inbound.Country)
		}
	}
	return vpnSources
}

func (i *HostInbound) countRequiredPortEachHost(inbounds []*entity.Inbound, VPNSourceDomains map[string][]string) map[string]uint32 {
	ret := map[string]uint32{}
	vpnSourceRequiredPortsCount := i.countInboundsByVPNSource(inbounds)

	for vpnSource, domains := range VPNSourceDomains {
		domainPortsCount := i.countDominPortsRequired(
			domains,
			vpnSourceRequiredPortsCount[vpnSource],
		)
		ret = i.mergeTwoDomainPortsCount(ret, domainPortsCount)
	}
	return ret
}

func (i *HostInbound) countDominPortsRequired(domains []string, requiredPortsCount uint32) map[string]uint32 {
	ret := map[string]uint32{}
	for _, domain := range domains {
		ret[domain] += requiredPortsCount
	}

	return ret
}

func (i *HostInbound) mergeTwoDomainPortsCount(a, b map[string]uint32) map[string]uint32 {
	ret := map[string]uint32{}
	for domain, count := range a {
		ret[domain] += count
	}

	for domain, count := range b {
		ret[domain] += count
	}

	return ret
}

func (i *HostInbound) countInboundsByVPNSource(inbounds []*entity.Inbound) map[string]uint32 {
	ret := map[string]uint32{}
	for _, inbound := range inbounds {
		ret[inbound.Country] += 1
	}
	return ret
}

func (i *HostInbound) getVPNSourceInboundDestination(
	VPNSourceDomains map[string][]string,
	hostPorts map[string][]string,
) map[string][]*Address {
	VPNSourceInboundDestination := map[string][]*Address{}
	for country, domains := range VPNSourceDomains {
		VPNSourceInboundDestination[country] = i.getAddressesByCountry(domains, hostPorts)
	}
	return VPNSourceInboundDestination
}

func (i *HostInbound) getAddressesByCountry(domains []string, hostPorts map[string][]string) []*Address {
	addresses := []*Address{}
	for _, domain := range domains {
		addresses = append(
			addresses,
			i.makeAddresses(domain, hostPorts[domain])...,
		)
	}
	return addresses
}

func (i *HostInbound) makeAddresses(domain string, ports []string) []*Address {
	ret := []*Address{}
	for _, port := range ports {
		ret = append(ret, &Address{Domain: domain, Port: port})
	}
	return ret
}

func (i *HostInbound) getDomainPorts(inbounds []*entity.Inbound) map[string][]string {
	domainPorts := make(map[string][]string)
	for _, inbound := range inbounds {
		domainPorts[inbound.Domain] = append(domainPorts[inbound.Domain], inbound.Port)
	}
	return domainPorts
}

func (i *HostInbound) getHostPortsMap(hostPortsFailures []*hostServiceDto.HostPortsFailed) map[string]map[string]struct{} {
	hostPortsMap := make(map[string]map[string]struct{})
	for _, hostPortFailed := range hostPortsFailures {
		hostPortsMap[hostPortFailed.Domain] = i.getPortFailedMap(hostPortFailed)
	}
	return hostPortsMap
}

func (i *HostInbound) getPortFailedMap(hostPortsFailed *hostServiceDto.HostPortsFailed) map[string]struct{} {
	portMap := make(map[string]struct{})
	for _, port := range hostPortsFailed.Ports {
		portMap[port] = struct{}{}
	}
	return portMap
}
