package worker

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mohamadrezamomeni/momo/contract/gogrpc/metric"
	"github.com/mohamadrezamomeni/momo/contract/gogrpc/port"
	"github.com/mohamadrezamomeni/momo/entity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	port.UnimplementedPortServer
	metric.UnimplementedMetricServer
	metricSvc MetricService
	portSvc   PortService
	address   string
}

type MetricService interface {
	GetMetric() (int, entity.HostStatus, error)
}

type PortService interface {
	GetAvailablePorts(uint32, []string) ([]string, error)
	OpenPorts([]string) []string
}

func New(metricSvc MetricService, portSvc PortService, metricConfig WorkerConfig) *Server {
	address := fmt.Sprintf("%s:%s", metricConfig.Address, metricConfig.Port)
	return &Server{
		UnimplementedMetricServer: metric.UnimplementedMetricServer{},
		UnimplementedPortServer:   port.UnimplementedPortServer{},
		metricSvc:                 metricSvc,
		portSvc:                   portSvc,
		address:                   address,
	}
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	metric.RegisterMetricServer(server, s)
	port.RegisterPortServer(server, s)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) GetMetric(ctx context.Context, req *metric.MetricRequest) (*metric.MetricResponse, error) {
	return &metric.MetricResponse{
		Rank:   2,
		Status: entity.HostStatusString(entity.High),
	}, nil
}

func (s *Server) GetAvailablePorts(ctx context.Context, req *port.PortAssignRequest) (*port.PortAssignResponse, error) {
	ports, err := s.portSvc.GetAvailablePorts(req.RequestNumberOfPort, req.PortsUsed)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "invalid input: %v", err)
	}
	return &port.PortAssignResponse{
		Ports: ports,
	}, nil
}
