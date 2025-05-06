package httpserver

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/auth"
	metricHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/metric"
	userHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/user"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	authSvc "github.com/mohamadrezamomeni/momo/service/auth"
	cryptService "github.com/mohamadrezamomeni/momo/service/crypt"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	authValidation "github.com/mohamadrezamomeni/momo/validator/auth"
	userValidation "github.com/mohamadrezamomeni/momo/validator/user"
)

type Server struct {
	Router        *echo.Echo
	config        *HTTPConfig
	metricHandler *metricHandler.Handler
	userHandler   *userHandler.Handler
	authHandler   *authHandler.Handler
}

func New(cfg *HTTPConfig,
	authSvc *authSvc.Auth,
	userSvc *userService.User,
	crypt *cryptService.Crypt,
	userValidation *userValidation.Validator,
	authValidator *authValidation.Validation,
) *Server {
	return &Server{
		Router:        echo.New(),
		config:        cfg,
		metricHandler: metricHandler.New(),
		userHandler:   userHandler.New(userSvc, userValidation, authSvc),
		authHandler:   authHandler.New(authSvc, authValidator),
	}
}

func (s *Server) Serve() {
	scope := "httpserver.serve"

	s.Router.Use(middleware.RequestID())
	s.Router.Use(middleware.Recover())

	api := s.Router.Group("/api/v1")

	s.metricHandler.SetRouter(api)
	s.userHandler.SetRouter(api)
	s.authHandler.SetRouter(api)

	address := fmt.Sprintf(":%s", s.config.Port)
	if err := s.Router.Start(address); err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Router.Shutdown(ctx)
}
