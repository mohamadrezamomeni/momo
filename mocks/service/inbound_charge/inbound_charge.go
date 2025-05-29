package inboundcharge

import "github.com/mohamadrezamomeni/momo/entity"

type InboundChargeMock struct{}

func New() *InboundChargeMock {
	return &InboundChargeMock{}
}

func (ic *InboundChargeMock) ChargeInbound(_ *entity.Charge) error {
	return nil
}
