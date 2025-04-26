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
	var cpu, memory, network int

	cpu, err := m.getCPUMetric()
	if err != nil {
		return 0, entity.Uknown, err
	}

	memory, err = m.getMemoryMetric()
	if err != nil {
		return 0, entity.Uknown, err
	}

	network, err = m.getNetworkMetric()
	if err != nil {
		return 0, entity.Uknown, err
	}

	totalMetric := m.getTotalMetric(cpu, memory, network)
	status := m.getStatus(totalMetric)
	return totalMetric, status, nil
}

func (m *Metric) getStatus(totalMetric int) entity.HostStatus {
	if totalMetric >= m.metricConfig.HighStatus {
		return entity.High
	} else if totalMetric >= m.metricConfig.MediumStatus {
		return entity.Medium
	} else if totalMetric >= m.metricConfig.LowStatus {
		return entity.Low
	}
	return entity.Deactive
}

func (m *Metric) getTotalMetric(cpu int, memory int, network int) int {
	return m.metricConfig.CPUWeight*cpu + m.metricConfig.MemoryWeight*memory + network*m.metricConfig.NetworkWeight
}

func (m *Metric) getCPUMetric() (int, error) {
	return 0, nil
}

func (m *Metric) getNetworkMetric() (int, error) {
	return 0, nil
}

func (m *Metric) getMemoryMetric() (int, error) {
	return 0, nil
}
