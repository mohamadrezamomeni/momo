package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	metricHandler "github.com/mohamadrezamomeni/momo/delivery/http_server/controller/metric"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type Server struct {
	router        *echo.Echo
	config        *HTTPConfig
	metricHandler *metricHandler.Handler
}

type UserService interface{}

func New(cfg *HTTPConfig) *Server {
	return &Server{
		router:        echo.New(),
		config:        cfg,
		metricHandler: metricHandler.New(),
	}
}

func (s *Server) Serve() {
	scope := "httpserver.serve"

	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Recover())

	api := s.router.Group("/api")

	s.metricHandler.SetRouter(api)

	address := fmt.Sprintf(":%s", s.config.Port)
	if err := s.router.Start(address); err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}
}
