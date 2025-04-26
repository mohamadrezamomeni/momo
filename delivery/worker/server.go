package worker

import (
	"context"
	"fmt"
	"log"
	"net"

	"momo/contract/gogrpc/metric"
	"momo/contract/gogrpc/port"
	"momo/entity"

	"google.golang.org/grpc"
)

type Server struct {
	port.UnimplementedPortServer
	metric.UnimplementedMetricServer
	metricSvc metricService
	portSvc   portService
	address   string
}

type metricService interface {
	GetMetric() (int, entity.HostStatus, error)
}

type portService interface {
	GetAvailablePort() (string, error)
}

func New(metricSvc metricService, portSvc portService, metricConfig WorkerConfig) *Server {
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

func (s *Server) GetAvailablePort(ctx context.Context, req *port.PortAssignRequest) (*port.PortAssignResponse, error) {
	p, err := s.portSvc.GetAvailablePort()
	if err != nil {
		return nil, err
	}
	return &port.PortAssignResponse{
		Port: p,
	}, nil
}
