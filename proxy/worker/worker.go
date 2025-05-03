package worker

import (
	"fmt"

	metric "github.com/mohamadrezamomeni/momo/contract/gogrpc/metric"
	port "github.com/mohamadrezamomeni/momo/contract/gogrpc/port"

	"google.golang.org/grpc"
)

type ProxyWorker struct {
	conn         *grpc.ClientConn
	portClient   port.PortClient
	metricClient metric.MetricClient
	address      string
}

func New(cfg *Config) (*ProxyWorker, error) {
	address := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return &ProxyWorker{}, err
	}

	return &ProxyWorker{
		portClient:   port.NewPortClient(conn),
		metricClient: metric.NewMetricClient(conn),
		conn:         conn,
		address:      address,
	}, nil
}

func (ps *ProxyWorker) Close() {
	ps.conn.Close()
}
