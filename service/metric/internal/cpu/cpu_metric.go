package cpu

func (c *CpuMetric) GetCpuFreePercentage() (float64, error) {
	total, idle, err := c.getData()
	if err != nil {
		return 0, err
	}
	return (float64(idle) / float64(total)) * 100, nil
}
