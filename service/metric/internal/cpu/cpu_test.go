package cpu

import (
	"testing"
)

func TestGetCpuMetric(t *testing.T) {
	cpuMetric, err := New()
	if err != nil {
		t.Fatalf("cpu metric isn't initiazed")
	}
	cpuFreeUsed := cpuMetric.GetCpuFreePercentage()

	if cpuFreeUsed < 0 {
		t.Fatal("error to get cpuFreeUsed")
	}
}
