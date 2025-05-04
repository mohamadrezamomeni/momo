package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	metricHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/metric"
	userHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/user"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	authSvc "github.com/mohamadrezamomeni/momo/service/auth"
	cryptService "github.com/mohamadrezamomeni/momo/service/crypt"
	userService "github.com/mohamadrezamomeni/momo/service/user"
	userValidation "github.com/mohamadrezamomeni/momo/validator/user"
)

type Server struct {
	router        *echo.Echo
	config        *HTTPConfig
	metricHandler *metricHandler.Handler
	userHandler   *userHandler.Handler
}

func New(cfg *HTTPConfig,
	authSvc *authSvc.Auth,
	userSvc *userService.User,
	crypt *cryptService.Crypt,
	userValidation *userValidation.Validator,
) *Server {
	return &Server{
		router:        echo.New(),
		config:        cfg,
		metricHandler: metricHandler.New(),
		userHandler:   userHandler.New(userSvc, userValidation, authSvc),
	}
}

func (s *Server) Serve() {
	scope := "httpserver.serve"

	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Recover())

	api := s.router.Group("/api/v1")

	s.metricHandler.SetRouter(api)
	s.userHandler.SetHandler(api)

	address := fmt.Sprintf(":%s", s.config.Port)
	if err := s.router.Start(address); err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}
}
