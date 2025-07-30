package charge

import (
	chargeRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/charge"
	"github.com/mohamadrezamomeni/momo/entity"
)

var (
	charge1 = &chargeRepositoryDto.CreateDto{
		Status:    entity.PendingStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Detail:    "hello",
		InboundID: "12",
		PackageID: "1",
	}

	charge2 = &chargeRepositoryDto.CreateDto{
		Status:    entity.PendingStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Detail:    "hi",
		InboundID: "12",
		PackageID: "3",
	}

	charge3 = &chargeRepositoryDto.CreateDto{
		Status:    entity.PendingStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Detail:    "hi",
		InboundID: "13",
		PackageID: "4",
	}

	charge4 = &chargeRepositoryDto.CreateDto{
		Status:    entity.ApprovedStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Detail:    "hello",
		InboundID: "12",
		PackageID: "1",
	}
	charge5 = &chargeRepositoryDto.CreateDto{
		Status:    entity.ApprovedStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Detail:    "hi",
		InboundID: "12",
		PackageID: "3",
	}

	charge6 = &chargeRepositoryDto.CreateDto{
		Status:    entity.ApprovedStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Detail:    "hi",
		InboundID: "15",
		PackageID: "3",
	}

	charge7 = &chargeRepositoryDto.CreateDto{
		Status:    entity.ApprovedStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Detail:    "hi",
		InboundID: "15",
		PackageID: "3",
	}

	charge8 = &chargeRepositoryDto.CreateDto{
		Status:    entity.ApprovedStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Detail:    "hi",
		InboundID: "15",
		PackageID: "3",
	}

	charge9 = &chargeRepositoryDto.CreateDto{
		Status:    entity.ApprovedStatusCharge,
		UserID:    "f47ac10b-58cc-4372-a567-0e02b2c3d477",
		Detail:    "hi",
		InboundID: "16",
		PackageID: "3",
	}
	charge10 = &chargeRepositoryDto.CreateDto{
		VPNType: entity.XRAY_VPN,
		Status:  entity.ApprovedStatusCharge,
		UserID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Detail:  "hello",
		Country: "1000",
	}
	charge11 = &chargeRepositoryDto.CreateDto{
		VPNType: entity.XRAY_VPN,
		Status:  entity.ApprovedStatusCharge,
		UserID:  "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Detail:  "hello",
		Country: "1000",
	}
)
