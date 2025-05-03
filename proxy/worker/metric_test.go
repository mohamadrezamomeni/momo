package worker

import (
	"testing"

	"github.com/mohamadrezamomeni/momo/entity"
)

func TestGetMetric(t *testing.T) {
	rank, status, err := pw.GetMetric()
	if err != nil {
		t.Fatalf("something has happend we couldn't get metric the problem was %v", err)
	}
	if rank < 0 {
		t.Error("rank is wrong.")
	}
	if status == entity.Uknown {
		t.Fatal("we got unkhon status")
	}
}
