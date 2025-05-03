package worker

import (
	"context"
	"time"

	"github.com/mohamadrezamomeni/momo/contract/gogrpc/metric"
	"github.com/mohamadrezamomeni/momo/entity"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (ps *ProxyWorker) GetMetric() (uint32, entity.HostStatus, error) {
	scope := "workerProxy.GetMetric"

	metricClient := metric.NewMetricClient(ps.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	metric, err := metricClient.GetMetric(ctx, &metric.MetricRequest{})
	if err != nil {
		return 0, entity.Uknown, momoError.Wrap(err).Scope(scope).Errorf("the address is %s", ps.address)
	}
	status, err := entity.MapHostStatusToEnum(metric.Status)
	if err != nil {
		return 0, entity.Uknown, momoError.Wrap(err).Scope(scope).Errorf("the address is %s", ps.address)
	}
	return metric.Rank, status, nil
}
