package inbound

import (
	"strconv"

	"github.com/mohamadrezamomeni/momo/entity"
)

type HostInbound struct {
	inboundRepo   InboundHostRepo
	vpnManagerSvc VPNManagerService
}

type AvailbleVPNAddress struct {
	VPNType entity.VPNType
	Domain  string
	Port    string
}

type VPNManagerService interface {
	GetAvailableVPNSourceDomains(vpnsSources []string, vpnTypes []entity.VPNType) ([]*entity.VPN, error)
}

type InboundHostRepo interface {
	FindInboundIsNotAssigned() ([]*entity.Inbound, error)
	GetListOfPortsByDomain() ([]struct {
		Domain string
		Ports  []string
	}, error)
	UpdateDomainPort(int, string, string) error
}

type Address struct {
	Domain string
	Port   string
}

func NewHostInbound(
	inboundRepo InboundHostRepo,
	VPNManagerSvc VPNManagerService,
) *HostInbound {
	return &HostInbound{
		inboundRepo:   inboundRepo,
		vpnManagerSvc: VPNManagerSvc,
	}
}

func (hi *HostInbound) AssignDomainToInbounds() {
	inbounds, err := hi.inboundRepo.FindInboundIsNotAssigned()
	if err != nil {
		return
	}

	portUserSummery, err := hi.summeryDomainPorts()
	if err != nil {
		return
	}
	vpnTypes := hi.getVPNTypes(inbounds)
	vpnSources := hi.getVPNSourcesFromInbounds(inbounds)

	vpns, err := hi.vpnManagerSvc.GetAvailableVPNSourceDomains(vpnSources, vpnTypes)
	if err != nil {
		return
	}
	availbleVPNsAddresses := hi.getAvailbleAddressesByVPNPortCountry(vpns, portUserSummery)

	hi.setAddresses(inbounds, availbleVPNsAddresses)
}

func (hi *HostInbound) setAddresses(inbounds []*entity.Inbound, availbleVPNsAddresses map[int]map[string][]*AvailbleVPNAddress) {
	indexes := make(map[entity.VPNType]map[string]int)
	for _, inbound := range inbounds {
		availableVPNPorts := availbleVPNsAddresses[inbound.VPNType][inbound.Country]
		if _, isExistVPNType := indexes[inbound.VPNType]; !isExistVPNType {
			indexes[inbound.VPNType] = make(map[string]int)
		}
		if _, isExistCountry := indexes[inbound.VPNType][inbound.Country]; !isExistCountry {
			indexes[inbound.VPNType][inbound.Country] = 0
		}
		i := indexes[inbound.VPNType][inbound.Country]
		if len(availableVPNPorts) > i {
			address := availbleVPNsAddresses[inbound.VPNType][inbound.Country][i]
			indexes[inbound.VPNType][inbound.Country] = i + 1
			hi.inboundRepo.UpdateDomainPort(inbound.ID, address.Domain, address.Port)
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

func (i *HostInbound) getVPNTypes(inbounds []*entity.Inbound) []entity.VPNType {
	res := make([]entity.VPNType, 0)
	seen := make(map[entity.VPNType]struct{})
	for _, inbound := range inbounds {
		if _, isExist := seen[inbound.VPNType]; !isExist {
			res = append(res, inbound.VPNType)
			seen[inbound.VPNType] = struct{}{}
		}
	}
	return res
}

func (i *HostInbound) getAvailbleAddressesByVPNPortCountry(vpns []*entity.VPN, domainPortSummery map[string][]string) map[entity.VPNType]map[string][]*AvailbleVPNAddress {
	res := make(map[entity.VPNType]map[string][]*AvailbleVPNAddress)
	for _, vpn := range vpns {
		if _, isExistCountry := res[vpn.VPNType][vpn.Country]; !isExistCountry {
			res[vpn.VPNType] = make(map[string][]*AvailbleVPNAddress, 0)
			res[vpn.VPNType][vpn.Country] = make([]*AvailbleVPNAddress, 0)
		}

		res[vpn.VPNType][vpn.Country] = append(
			res[vpn.VPNType][vpn.Country],
			i.getAvailblePortsByVPN(vpn, domainPortSummery)...,
		)
	}
	return res
}

func (i *HostInbound) getAvailblePortsByVPN(vpn *entity.VPN, domainPortSummery map[string][]string) []*AvailbleVPNAddress {
	domain := vpn.Domain
	portsUsed := domainPortSummery[domain]
	portUsedMap := i.getPortUsedMap(portsUsed)
	res := make([]*AvailbleVPNAddress, 0)

	for p := vpn.StartPort; p < vpn.EndPort+1; p++ {
		pStr := strconv.Itoa(p)
		if _, isExist := portUsedMap[pStr]; !isExist {
			res = append(res, &AvailbleVPNAddress{
				Domain:  domain,
				VPNType: vpn.VPNType,
				Port:    strconv.Itoa(p),
			})
		}
	}
	return res
}

func (i *HostInbound) getPortUsedMap(portsUsed []string) map[string]struct{} {
	portUsedMap := make(map[string]struct{})
	for _, p := range portsUsed {
		portUsedMap[p] = struct{}{}
	}
	return portUsedMap
}
