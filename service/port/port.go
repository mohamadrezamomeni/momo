package port

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"momo/delivery/worker"
	momoError "momo/pkg/error"
)

type Port struct {
	mu        sync.Mutex
	startPort int
	endPort   int
}

func New(cfg *worker.PortAssignment) *Port {
	return &Port{
		startPort: cfg.StartPort,
		endPort:   cfg.EndPort,
	}
}

func (p *Port) GetAvailablePort() (string, error) {
	for i := p.startPort; i < p.endPort+1; i++ {
		if p.isPortAvailable(i) {
			return strconv.Itoa(i), nil
		}
	}
	return "", momoError.Error("we couldn't find available port")
}

func (p *Port) isPortAvailable(port int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	defer listen.Close()

	return true
}
