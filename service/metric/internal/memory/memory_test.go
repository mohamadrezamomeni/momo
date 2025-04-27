package memory

import "testing"

func TestIntializeMemoryMetric(t *testing.T) {
	mm, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if !(mm.MemAvailable >= 0 && mm.MemFree >= 0 && mm.MemTotal >= 0) {
		t.Fatalf("something went wrong to give information")
	}
}
