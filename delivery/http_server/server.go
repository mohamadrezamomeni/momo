package httpserver

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/auth"
	hostHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/host"
	inboundHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/inbound"
	metricHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/metric"
	userHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/user"
	vpnHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/vpn"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	authSvc "github.com/mohamadrezamomeni/momo/service/auth"
	cryptService "github.com/mohamadrezamomeni/momo/service/crypt"
	hostService "github.com/mohamadrezamomeni/momo/service/host"
	inboundService "github.com/mohamadrezamomeni/momo/service/inbound"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	vpnService "github.com/mohamadrezamomeni/momo/service/vpn_manager"
	authValidation "github.com/mohamadrezamomeni/momo/validator/auth"
	hostValidation "github.com/mohamadrezamomeni/momo/validator/host"
	inboundValidation "github.com/mohamadrezamomeni/momo/validator/inbound"
	userValidation "github.com/mohamadrezamomeni/momo/validator/user"
	vpnValidation "github.com/mohamadrezamomeni/momo/validator/vpn"
)

type Server struct {
	router         *echo.Echo
	config         *HTTPConfig
	metricHandler  *metricHandler.Handler
	userHandler    *userHandler.Handler
	authHandler    *authHandler.Handler
	hostHandler    *hostHandler.Handler
	vpnHandler     *vpnHandler.Handler
	inboundHandler *inboundHandler.Handler
}

func New(cfg *HTTPConfig,
	authSvc *authSvc.Auth,
	userSvc *userService.User,
	crypt *cryptService.Crypt,
	hostSvc *hostService.Host,
	vpnSvc *vpnService.VPNService,
	inboundSvc *inboundService.Inbound,
	userValidation *userValidation.Validator,
	authValidator *authValidation.Validation,
	hostValidator *hostValidation.Validator,
	vpnValidator *vpnValidation.Validator,
	inboundValidator *inboundValidation.Validator,
) *Server {
	return &Server{
		router:         echo.New(),
		config:         cfg,
		metricHandler:  metricHandler.New(),
		userHandler:    userHandler.New(userSvc, userValidation, authSvc),
		authHandler:    authHandler.New(authSvc, authValidator),
		hostHandler:    hostHandler.New(hostSvc, authSvc, hostValidator),
		vpnHandler:     vpnHandler.New(vpnSvc, vpnValidator, authSvc),
		inboundHandler: inboundHandler.New(inboundSvc, inboundValidator, authSvc),
	}
}

func (s *Server) Serve() {
	scope := "httpserver.serve"

	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Recover())

	api := s.router.Group("/api/v1")

	s.metricHandler.SetRouter(api)
	s.userHandler.SetRouter(api)
	s.authHandler.SetRouter(api)
	s.hostHandler.SetRouter(api)
	s.vpnHandler.SetRouter(api)
	s.inboundHandler.SetRouter(api)

	address := fmt.Sprintf(":%s", s.config.Port)
	if err := s.router.Start(address); err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
