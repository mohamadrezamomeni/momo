package worker

import (
	"context"
	"time"

	"momo/contract/gogrpc/metric"
	"momo/entity"
	momoError "momo/pkg/error"
)

func (ps *ProxyWorker) GetMetric() (uint32, entity.HostStatus, error) {
	metricClient := metric.NewMetricClient(ps.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	metric, err := metricClient.GetMetric(ctx, &metric.MetricRequest{})
	if err != nil {
		return 0, entity.Uknown, momoError.Errorf("something wrong has happend to get metric from %s, the error was %v", ps.address, err)
	}
	status, err := entity.MapHostStatusToEnum(metric.Status)
	if err != nil {
		return 0, entity.Uknown, err
	}
	return metric.Rank, status, nil
}
