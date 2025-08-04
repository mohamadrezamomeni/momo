package inbound

import (
	"fmt"
	"strconv"

	inboundServiceDto "github.com/mohamadrezamomeni/momo/dto/service/inbound"
	"github.com/mohamadrezamomeni/momo/entity"
)

type MockInboundService struct {
	inbounds []*entity.Inbound
	idx      int
}

func New() *MockInboundService {
	return &MockInboundService{
		inbounds: make([]*entity.Inbound, 0),
		idx:      0,
	}
}

func (i *MockInboundService) Create(inpt *inboundServiceDto.CreateInbound) (*entity.Inbound, error) {
	inbound := &entity.Inbound{
		ID:         strconv.Itoa(i.idx),
		Domain:     "",
		Protocol:   "vmess",
		UserID:     inpt.UserID,
		IsActive:   false,
		IsNotified: false,
		IsAssigned: false,
		Port:       "",
		Start:      inpt.Start,
		End:        inpt.End,
		VPNType:    inpt.VPNType,
	}
	i.idx += 1
	i.inbounds = append(i.inbounds, inbound)
	return inbound, nil
}

func (i *MockInboundService) FindInboundByID(id string) (*entity.Inbound, error) {
	for _, inbound := range i.inbounds {
		if inbound.ID == id {
			return inbound, nil
		}
	}
	return nil, fmt.Errorf("")
}

func (i *MockInboundService) DeletedAll() {
	i.inbounds = make([]*entity.Inbound, 0)
}
