package worker

import (
	"context"
	"time"

	"momo/contract/gogrpc/metric"
	momoError "momo/pkg/error"
)

func (ps *ProxyWorker) GetMetric() (uint32, string, error) {
	metricClient := metric.NewMetricClient(ps.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	metric, err := metricClient.GetMetric(ctx, &metric.MetricRequest{})
	if err != nil {
		return 0, "", momoError.Errorf("something wrong has happend to get metric from %s, the error was %v", ps.address, err)
	}
	return metric.Rank, metric.Status, nil
}
