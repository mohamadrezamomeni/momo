package httpserver

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/auth"
	chargeHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/charge"
	hostHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/host"
	inboundHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/inbound"
	metricHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/metric"
	tierHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/tier"
	userHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/user"
	userTierHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/user_tier"
	vpnHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/vpn"
	vpnPackageHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/vpn_package"
	vpnSourceHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/vpn_source"
	momoLog "github.com/mohamadrezamomeni/momo/pkg/log"
	authSvc "github.com/mohamadrezamomeni/momo/service/auth"
	chargeService "github.com/mohamadrezamomeni/momo/service/charge"
	cryptService "github.com/mohamadrezamomeni/momo/service/crypt"
	hostService "github.com/mohamadrezamomeni/momo/service/host"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	inboundChargeService "github.com/mohamadrezamomeni/momo/service/inbound_charge"
	tierService "github.com/mohamadrezamomeni/momo/service/tier"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	userTierService "github.com/mohamadrezamomeni/momo/service/user_tiers"
	vpnService "github.com/mohamadrezamomeni/momo/service/vpn_manager"
	vpnPackageService "github.com/mohamadrezamomeni/momo/service/vpn_package"
	vpnSourceService "github.com/mohamadrezamomeni/momo/service/vpn_source"
	authValidation "github.com/mohamadrezamomeni/momo/validator/auth"
	chargeValidation "github.com/mohamadrezamomeni/momo/validator/charge"
	hostValidation "github.com/mohamadrezamomeni/momo/validator/host"
	inboundValidation "github.com/mohamadrezamomeni/momo/validator/inbound"
	userValidation "github.com/mohamadrezamomeni/momo/validator/user"
	vpnValidation "github.com/mohamadrezamomeni/momo/validator/vpn"
	vpnPackageValidation "github.com/mohamadrezamomeni/momo/validator/vpn_package"
	vpnSourceValidation "github.com/mohamadrezamomeni/momo/validator/vpn_source"
)

type Server struct {
	router            *echo.Echo
	config            *HTTPConfig
	metricHandler     *metricHandler.Handler
	userHandler       *userHandler.Handler
	authHandler       *authHandler.Handler
	hostHandler       *hostHandler.Handler
	vpnHandler        *vpnHandler.Handler
	inboundHandler    *inboundHandler.Handler
	chargeHandler     *chargeHandler.Handler
	vpnSourceHandler  *vpnSourceHandler.Handler
	vpnPackageHandler *vpnPackageHandler.Handler
	tierHandler       *tierHandler.Handler
	userTierHandler   *userTierHandler.Handler
}

func New(cfg *HTTPConfig,
	authSvc *authSvc.Auth,
	userSvc *userService.User,
	crypt *cryptService.Crypt,
	hostSvc *hostService.Host,
	vpnSvc *vpnService.VPNService,
	inboundSvc *inboundService.Inbound,
	inboundChargeSvc *inboundChargeService.InboundCharge,
	chargeSvc *chargeService.Charge,
	vpnPackageSvc *vpnPackageService.VPNPackage,
	vpnSourceSvc *vpnSourceService.VPNSource,
	tierSvc *tierService.Tier,
	userTierSvc *userTierService.UserTiers,
	userValidation *userValidation.Validator,
	authValidator *authValidation.Validation,
	hostValidator *hostValidation.Validator,
	vpnValidator *vpnValidation.Validator,
	inboundValidator *inboundValidation.Validator,
	vpnSourceValidator *vpnSourceValidation.Validator,
	vpnPackageValidator *vpnPackageValidation.Validator,
	chargeValidator *chargeValidation.Validator,
) *Server {
	return &Server{
		router:            echo.New(),
		config:            cfg,
		metricHandler:     metricHandler.New(),
		userHandler:       userHandler.New(userSvc, userValidation, authSvc),
		authHandler:       authHandler.New(authSvc, authValidator),
		hostHandler:       hostHandler.New(hostSvc, authSvc, hostValidator),
		vpnHandler:        vpnHandler.New(vpnSvc, vpnValidator, authSvc),
		inboundHandler:    inboundHandler.New(inboundSvc, inboundValidator, authSvc),
		chargeHandler:     chargeHandler.New(chargeSvc, authSvc, chargeValidator, inboundChargeSvc),
		vpnSourceHandler:  vpnSourceHandler.New(vpnSourceSvc, vpnSourceValidator, authSvc),
		vpnPackageHandler: vpnPackageHandler.New(vpnPackageSvc, vpnPackageValidator, authSvc),
		tierHandler:       tierHandler.New(tierSvc, authSvc),
		userTierHandler:   userTierHandler.New(userTierSvc, authSvc),
	}
}

func (s *Server) Serve() {
	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Recover())

	api := s.router.Group("/api/v1")

	s.metricHandler.SetRouter(api)
	s.userHandler.SetRouter(api)
	s.authHandler.SetRouter(api)
	s.hostHandler.SetRouter(api)
	s.vpnHandler.SetRouter(api)
	s.inboundHandler.SetRouter(api)
	s.chargeHandler.SetRouter(api)
	s.vpnSourceHandler.SetRouter(api)
	s.vpnPackageHandler.SetRouter(api)
	s.tierHandler.SetRouter(api)
	s.userTierHandler.SetRouter(api)

	address := fmt.Sprintf(":%s", s.config.Port)
	if err := s.router.Start(address); err != nil {
		momoLog.Info(err.Error())
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
