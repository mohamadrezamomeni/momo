package inbound

import (
	"github.com/mohamadrezamomeni/momo/entity"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type HostInbound struct {
	inboundRepo InboundHostRepo
	hostService HostService
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
	ResolveHostPortPair(map[string][]string, int) ([][2]string, error)
}

func NewHostInbound(
	inboundRepo InboundHostRepo,
	hostService HostService,
) *HostInbound {
	return &HostInbound{
		inboundRepo: inboundRepo,
		hostService: hostService,
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

	hostPortPairs, err := i.hostService.ResolveHostPortPair(portUserSummery, len(inbounds))
	if err != nil {
		return
	}

	for j := 0; j < utils.Min(len(inbounds), len(hostPortPairs)); j++ {
		hostPort := hostPortPairs[j]
		inbound := inbounds[j]
		host, port := hostPort[0], hostPort[1]
		i.inboundRepo.UpdateDomainPort(inbound.ID, host, port)
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
