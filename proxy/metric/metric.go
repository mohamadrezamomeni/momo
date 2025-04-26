package metric

import (
	"context"
	"fmt"
	"time"

	"momo/contract/gogrpc/metric"

	momoError "momo/pkg/error"

	"google.golang.org/grpc"
)

type ProxyMetric struct {
	conn    *grpc.ClientConn
	address string
}

func New(cfg *Config) (*ProxyMetric, error) {
	address := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return &ProxyMetric{}, err
	}

	return &ProxyMetric{
		conn:    conn,
		address: address,
	}, nil
}

func (ps *ProxyMetric) GetMetric() (uint32, string, error) {
	metricClient := metric.NewMetricClient(ps.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	metric, err := metricClient.GetMetric(ctx, &metric.MetricRequest{})
	if err != nil {
		return 0, "", momoError.Errorf("something wrong has happend to get metric from %s, the error was %v", ps.address, err)
	}
	return metric.Rank, metric.Status, nil
}

func (ps *ProxyMetric) Close() {
	ps.conn.Close()
}
