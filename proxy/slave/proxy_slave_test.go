package slaves

import (
	"os"
	"testing"

	"momo/entity"
)

var ps *ProxySlave

func TestMain(m *testing.M) {
	ps, _ = New(&Config{Address: "localhost", Port: "666"})

	code := m.Run()
	ps.Close()
	os.Exit(code)
}

func TestGetMetric(t *testing.T) {
	rank, status, err := ps.GetMetric()
	if err != nil {
		t.Errorf("something has happend we couldn't get metric the problem was %v", err)
	}
	if rank < 0 {
		t.Error("rank is wrong.")
	}
	hostStatus, _ := entity.MapHostStatusToEnum(status)
	if hostStatus == entity.Uknown {
		t.Fatal("we got unkhon status")
	}
}
