package worker

import (
	"os"
	"testing"
)

var pw *ProxyWorker

func TestMain(m *testing.M) {
	pw, _ = New(&Config{Address: "localhost", Port: "666"})

	code := m.Run()
	pw.Close()
	os.Exit(code)
}
