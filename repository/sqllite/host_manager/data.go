package hostmanager

import (
	hostmanagerDto "momo/dto/repository/host_manager"
	"momo/entity"
)

var (
	hostExample1 = &hostmanagerDto.AddHost{
		Domain:         "google.com",
		Port:           "62789",
		Status:         entity.Deactive,
		StartRangePort: 1000,
		EndRangePort:   2000,
	}
	hostExample2 = &hostmanagerDto.AddHost{
		Domain:         "yahoo.com",
		Port:           "62780",
		Status:         entity.High,
		StartRangePort: 2500,
		EndRangePort:   3000,
	}
	hostExample3 = &hostmanagerDto.AddHost{
		Domain:         "facebook.com",
		Port:           "62780",
		Status:         entity.Deactive,
		StartRangePort: 2000,
		EndRangePort:   5000,
	}
	hostExample4 = &hostmanagerDto.AddHost{
		Domain:         "twitter.com",
		Port:           "62780",
		Status:         entity.Medium,
		StartRangePort: 1000,
		EndRangePort:   2000,
	}
	hostExample5 = &hostmanagerDto.AddHost{
		Domain:         "github.com",
		Port:           "62780",
		Status:         entity.Low,
		StartRangePort: 1000,
		EndRangePort:   2000,
	}
	hostExample6 = &hostmanagerDto.AddHost{
		Domain:         "gitlab.com",
		Port:           "62780",
		Status:         entity.High,
		StartRangePort: 1000,
		EndRangePort:   2000,
	}
)
