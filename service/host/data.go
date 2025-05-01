package host

import (
	hostRepositoryDto "momo/dto/repository/host_manager"
	"momo/entity"
)

var (
	host1 = &hostRepositoryDto.AddHost{
		Domain: "twitter.com",
		Port:   "222",
		Status: entity.High,
		Rank:   10,
	}
	host2 = &hostRepositoryDto.AddHost{
		Domain: "facebook.com",
		Port:   "222",
		Status: entity.High,
		Rank:   10,
	}
	host3 = &hostRepositoryDto.AddHost{
		Domain: "google.com",
		Port:   "222",
		Status: entity.High,
		Rank:   10,
	}

	host4 = &hostRepositoryDto.AddHost{
		Domain: "google.com",
		Port:   "222",
		Status: entity.High,
		Rank:   7,
	}

	host5 = &hostRepositoryDto.AddHost{
		Domain: "google.com",
		Port:   "222",
		Status: entity.Low,
		Rank:   5,
	}

	host6 = &hostRepositoryDto.AddHost{
		Domain: "google.com",
		Port:   "222",
		Status: entity.Uknown,
		Rank:   0,
	}
)
