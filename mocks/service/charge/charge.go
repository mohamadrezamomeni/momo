package charge

import "github.com/mohamadrezamomeni/momo/entity"

type ChargeMock struct{}

func New() *ChargeMock {
	return &ChargeMock{}
}

func (c *ChargeMock) FindAvailbleCharge(_ string) (*entity.Charge, error) {
	return &entity.Charge{
		Detail:       "test",
		PackageID:    "11",
		InboundID:    "11",
		Status:       entity.ApprovedStatusCharge,
		UserID:       "12345678",
		AdminComment: "1234567",
		ID:           "12",
	}, nil
}
