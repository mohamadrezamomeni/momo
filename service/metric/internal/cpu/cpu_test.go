package cpu

import (
	"path/filepath"
	"testing"

	"momo/pkg/utils"
)

func TestGetCpuMetric(t *testing.T) {
	cpuMetric := New()
	root, err := utils.GetRootOfProject()
	if err != nil {
		t.Fatalf("cpu metric isn't initiazed")
	}
	cpuMetric.statFilePath = filepath.Join(root, "stat.test")
	cpuFreeUsed, err := cpuMetric.GetCpuFreePercentage()
	if err != nil {
		t.Fatal(err)
	}
	if cpuFreeUsed < 0 {
		t.Fatal("error to get cpuFreeUsed")
	}
}
