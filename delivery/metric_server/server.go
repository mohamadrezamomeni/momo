package metricserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"momo/contract/gogrpc/metric"
	"momo/entity"

	"google.golang.org/grpc"
)

type Server struct {
	metric.UnimplementedMetricServer
	svc     metricService
	address string
}

type metricService interface {
	GetMetric() (int, entity.HostStatus, error)
}

func New(metricSvc metricService, metricConfig MetricConfig) *Server {
	address := fmt.Sprintf("%s:%s", metricConfig.Address, metricConfig.Port)
	return &Server{
		UnimplementedMetricServer: metric.UnimplementedMetricServer{},
		svc:                       metricSvc,
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
