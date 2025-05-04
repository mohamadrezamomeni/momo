package service

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"

	hostManagerSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/host_manager"
	inboundSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/inbound"
	userSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/user"
	vpnSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_manager"

	cryptService "github.com/mohamadrezamomeni/momo/service/crypt"
	hostService "github.com/mohamadrezamomeni/momo/service/host"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	vpnService "github.com/mohamadrezamomeni/momo/service/vpn_manager"
)

func GetServices(cfg *sqllite.DBConfig) (
	*hostService.Host,
	*vpnService.VPNService,
	*userService.User,
	*inboundService.Inbound,
) {
	db := sqllite.New(cfg)
	userRepo := userSqlite.New(db)
	inboundRepo := inboundSqlite.New(db)
	hostRepo := hostManagerSqlite.New(db)
	vpnRepo := vpnSqlite.New(db)

	hostSvc := hostService.New(hostRepo, adapter.AdaptedWorkerFactory)
	vpnSvc := vpnService.New(vpnRepo, adapter.AdaptedVPNProxyFactory)
	userSvc := userService.New(userRepo, cryptService.New())

	inbouncSvc := inboundService.New(inboundRepo, vpnSvc, userSvc, hostSvc)
	return hostSvc, vpnSvc, userSvc, inbouncSvc
}
