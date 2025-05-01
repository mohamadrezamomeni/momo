package memory

import (
	"path/filepath"
	"testing"

	"momo/pkg/utils"
)

func TestIntializeMemoryMetric(t *testing.T) {
	mm := New()

	root, _ := utils.GetRootOfProject()
	mm.memInfoPath = filepath.Join(root, "meminfo.test")
	memTotal, memFree, memAvailable, err := mm.GetData()
	if err != nil {
		t.Fatal(err)
	}
	if !(memTotal >= 0 && memFree >= 0 && memAvailable >= 0) {
		t.Fatalf("something went wrong to give information")
	}
}
