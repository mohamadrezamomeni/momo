package cpu

func (m *CpuMetric) GetCpuFreePercentage() uint64 {
	return m.idle / m.total * 100
}
