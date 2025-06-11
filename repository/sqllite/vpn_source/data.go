package vpnsource

import vpnSourceRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"

var (
	vpnsource1 = &vpnSourceRepositoryDto.CreateVPNSourceDto{
		Country: "us",
		English: "us",
	}
	vpnsource2 = &vpnSourceRepositoryDto.CreateVPNSourceDto{
		Country: "uk",
		English: "uk",
	}
	vpnsource3 = &vpnSourceRepositoryDto.CreateVPNSourceDto{
		Country: "china",
		English: "china",
	}
)
