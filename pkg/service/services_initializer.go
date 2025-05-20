package service

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"

	eventSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/event"
	hostManagerSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/host_manager"
	inboundSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/inbound"
	userSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/user"
	vpnSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_manager"
	vpnPackageSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_package"

	config "github.com/mohamadrezamomeni/momo/pkg/config"
	authService "github.com/mohamadrezamomeni/momo/service/auth"
	cryptService "github.com/mohamadrezamomeni/momo/service/crypt"
	eventService "github.com/mohamadrezamomeni/momo/service/event"
	hostService "github.com/mohamadrezamomeni/momo/service/host"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	vpnService "github.com/mohamadrezamomeni/momo/service/vpn_manager"
	vpnPackageService "github.com/mohamadrezamomeni/momo/service/vpn_package"
)

func GetServices(cfg *config.Config) (
	*hostService.Host,
	*vpnService.VPNService,
	*userService.User,
	*inboundService.Inbound,
	*authService.Auth,
	*cryptService.Crypt,
	*vpnPackageService.VPNPackage,
	*eventService.Event,
) {
	db := sqllite.New(&cfg.DB)
	userRepo := userSqlite.New(db)
	inboundRepo := inboundSqlite.New(db)
	hostRepo := hostManagerSqlite.New(db)
	vpnRepo := vpnSqlite.New(db)
	vpnPackageRepo := vpnPackageSqlite.New(db)
	eventRepo := eventSqlite.New(db)

	hostSvc := hostService.New(hostRepo, adapter.AdaptedWorkerFactory)
	vpnSvc := vpnService.New(vpnRepo, adapter.AdaptedVPNProxyFactory)
	cryptSvc := cryptService.New(&cfg.CryptConfig)
	userSvc := userService.New(userRepo, cryptSvc)
	authSvc := authService.New(userSvc, cryptSvc, &cfg.AuthConfig)
	vpnPackageSvc := vpnPackageService.New(vpnPackageRepo)
	eventSvc := eventService.New(eventRepo)

	inbouncSvc := inboundService.New(inboundRepo, vpnSvc, userSvc, hostSvc)
	return hostSvc, vpnSvc, userSvc, inbouncSvc, authSvc, cryptSvc, vpnPackageSvc, eventSvc
}
