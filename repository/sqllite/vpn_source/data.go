package vpnsource

import vpnSourceRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"

var (
	vpnsource1 = &vpnSourceRepositoryDto.CreateVPNSourceDto{
		Title:   "us",
		English: "us",
	}
	vpnsource2 = &vpnSourceRepositoryDto.CreateVPNSourceDto{
		Title:   "uk",
		English: "uk",
	}
	vpnsource3 = &vpnSourceRepositoryDto.CreateVPNSourceDto{
		Title:   "china",
		English: "china",
	}
)
