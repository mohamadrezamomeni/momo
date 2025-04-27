package metric

import (
	metricServer "momo/delivery/worker"
	"momo/entity"
)

type Metric struct {
	metricConfig *metricServer.MetricConfig
}

func New(metricConfig *metricServer.MetricConfig) *Metric {
	return &Metric{
		metricConfig: metricConfig,
	}
}

func (m *Metric) GetMetric() (int, entity.HostStatus, error) {
	return 10, entity.High, nil
}
