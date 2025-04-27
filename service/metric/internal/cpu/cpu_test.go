package cpu

import "testing"

func TestNewCpuMetric(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatalf("error to initialize cpu Metric the problem was %v", err)
	}
}

func TestGetCpuMetric(t *testing.T) {
	cpuMetric, _ := New()

	cpuFreeUsed := cpuMetric.GetCpuFreePercentage()

	if cpuFreeUsed < 0 {
		t.Fatal("error to get cpuFreeUsed")
	}
}
