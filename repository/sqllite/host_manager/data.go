package hostmanager

import (
	hostmanagerDto "github.com/mohamadrezamomeni/momo/dto/repository/host_manager"
	"github.com/mohamadrezamomeni/momo/entity"
)

var (
	hostExample1 = &hostmanagerDto.AddHost{
		Domain: "google.com",
		Port:   "62789",
		Status: entity.Deactive,
	}
	hostExample2 = &hostmanagerDto.AddHost{
		Domain: "yahoo.com",
		Port:   "62780",
		Status: entity.High,
	}
	hostExample3 = &hostmanagerDto.AddHost{
		Domain: "facebook.com",
		Port:   "62780",
		Status: entity.Deactive,
	}
	hostExample4 = &hostmanagerDto.AddHost{
		Domain: "twitter.com",
		Port:   "62780",
		Status: entity.Medium,
	}
	hostExample5 = &hostmanagerDto.AddHost{
		Domain: "github.com",
		Port:   "62780",
		Status: entity.Low,
	}
	hostExample6 = &hostmanagerDto.AddHost{
		Domain: "gitlab.com",
		Port:   "62780",
		Status: entity.High,
	}
)
