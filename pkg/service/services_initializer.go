package service

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	"github.com/mohamadrezamomeni/momo/repository/sqllite"

	chargeSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/charge"
	eventSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/event"
	hostManagerSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/host_manager"
	inboundSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/inbound"
	inboundChargeSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/inbound_charge"
	tierSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/tier"
	userSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/user"
	userTierSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/user_tier"
	vpnSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_manager"
	vpnPackageSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_package"
	vpnSourceSqlite "github.com/mohamadrezamomeni/momo/repository/sqllite/vpn_source"

	config "github.com/mohamadrezamomeni/momo/pkg/config"
	authService "github.com/mohamadrezamomeni/momo/service/auth"
	chargeService "github.com/mohamadrezamomeni/momo/service/charge"
	cryptService "github.com/mohamadrezamomeni/momo/service/crypt"
	eventService "github.com/mohamadrezamomeni/momo/service/event"
	hostService "github.com/mohamadrezamomeni/momo/service/host"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	inboundChargeService "github.com/mohamadrezamomeni/momo/service/inbound_charge"
	tierService "github.com/mohamadrezamomeni/momo/service/tier"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	userTiersService "github.com/mohamadrezamomeni/momo/service/user_tiers"
	vpnService "github.com/mohamadrezamomeni/momo/service/vpn_manager"
	vpnPackageService "github.com/mohamadrezamomeni/momo/service/vpn_package"
	vpnSourceService "github.com/mohamadrezamomeni/momo/service/vpn_source"
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
	*chargeService.Charge,
	*inboundService.HealingUpInbound,
	*inboundService.HostInbound,
	*inboundService.InboundTraffic,
	*vpnSourceService.VPNSource,
	*inboundChargeService.InboundCharge,
	*tierService.Tier,
	*userTiersService.UserTiers,
) {
	db := sqllite.New(&cfg.DB)
	userRepo := userSqlite.New(db)
	inboundRepo := inboundSqlite.New(db)
	hostRepo := hostManagerSqlite.New(db)
	vpnRepo := vpnSqlite.New(db)
	vpnPackageRepo := vpnPackageSqlite.New(db)
	eventRepo := eventSqlite.New(db)
	chargeRepo := chargeSqlite.New(db)
	inboundChargeRepo := inboundChargeSqlite.New(db)
	vpnSourceRepo := vpnSourceSqlite.New(db)
	tierRepo := tierSqlite.New(db)
	userTierRepo := userTierSqlite.New(db)

	hostSvc := hostService.New(hostRepo, adapter.AdaptedWorkerFactory)
	vpnSvc := vpnService.New(vpnRepo, adapter.AdaptedVPNProxyFactory)
	cryptSvc := cryptService.New(&cfg.CryptConfig)

	userTiersSvc := userTiersService.New(userTierRepo)
	tierSvc := tierService.New(tierRepo)
	eventSvc := eventService.New(eventRepo)
	userSvc := userService.New(userRepo, cryptSvc, eventSvc)
	authSvc := authService.New(userSvc, cryptSvc, &cfg.AuthConfig)
	vpnPackageSvc := vpnPackageService.New(vpnPackageRepo)
	chargeSvc := chargeService.New(eventSvc, chargeRepo)
	inboundSvc := inboundService.New(inboundRepo)
	inboundChargeSvc := inboundChargeService.New(inboundChargeRepo, vpnPackageSvc, inboundRepo, chargeRepo, inboundSvc, eventSvc, chargeSvc)
	healingUpInboundSvc := inboundService.NewHealingUpInbound(inboundRepo, vpnSvc, inboundChargeSvc, chargeSvc, userSvc)
	hostInboundSvc := inboundService.NewHostInbound(inboundRepo, vpnSvc)
	inboundTrafficSvc := inboundService.NewInboundTraffic(inboundRepo, vpnSvc, userSvc)
	vpnSourceSvc := vpnSourceService.New(vpnSourceRepo, vpnSvc)

	return hostSvc,
		vpnSvc,
		userSvc,
		inboundSvc,
		authSvc,
		cryptSvc,
		vpnPackageSvc,
		eventSvc,
		chargeSvc,
		healingUpInboundSvc,
		hostInboundSvc,
		inboundTrafficSvc,
		vpnSourceSvc,
		inboundChargeSvc,
		tierSvc,
		userTiersSvc
}
