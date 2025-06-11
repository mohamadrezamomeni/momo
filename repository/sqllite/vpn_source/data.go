package vpnsource

import vpnSourceRepositoryDto "github.com/mohamadrezamomeni/momo/dto/repository/vpn_source"

var (
	vpnsource1 = &vpnSourceRepositoryDto.UpsertVPNSourceDto{
		English: "us",
	}
	vpnsource2 = &vpnSourceRepositoryDto.UpsertVPNSourceDto{
		English: "uk",
	}
	vpnsource3 = &vpnSourceRepositoryDto.UpsertVPNSourceDto{
		English: "china",
	}
)
