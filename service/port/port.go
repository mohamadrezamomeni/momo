package port

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	grpcWorker "github.com/mohamadrezamomeni/momo/delivery/grpc_worker"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type Port struct {
	mu        sync.Mutex
	startPort int
	endPort   int
}

func New(cfg *grpcWorker.PortAssignment) *Port {
	return &Port{
		startPort: cfg.StartPort,
		endPort:   cfg.EndPort,
	}
}

func (p *Port) GetAvailablePorts(portNeededCount uint32, portsUsed []string) ([]string, error) {
	mapPort := p.makeMapPorts(portsUsed)

	availblePorts := []string{}

	for i := p.startPort; i < p.endPort+1 && len(availblePorts) < int(portNeededCount); i++ {
		curPort := strconv.Itoa(i)
		p.mu.Lock()
		if _, ok := mapPort[curPort]; !ok && p.isPortAvailable(curPort) {
			availblePorts = append(availblePorts, curPort)
			p.store(curPort)
		}
		p.mu.Unlock()
	}

	return availblePorts, nil
}

func (p *Port) makeMapPorts(portsUsed []string) map[string]struct{} {
	hashMap := map[string]struct{}{}

	for _, port := range portsUsed {
		hashMap[port] = struct{}{}
	}

	return hashMap
}

func (p *Port) isPortAvailable(port string) bool {
	if p.isPortReserverd(port) {
		return false
	}
	if p.isPortBusy(port) {
		return false
	}
	return true
}

func (p *Port) isPortBusy(port string) bool {
	addr := fmt.Sprintf("127.0.0.1:%s", port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	defer listen.Close()
	return false
}

func (p *Port) isPortReserverd(port string) bool {
	oldestTime := p.findOldestRecordbyPort(port)
	if oldestTime.IsZero() {
		return false
	}

	existedPortAfterMinutes := oldestTime.Add(30 * time.Minute)

	now := time.Now()

	if existedPortAfterMinutes.Before(now) {
		return false
	}
	return true
}

func (p *Port) findOldestRecordbyPort(port string) time.Time {
	path := p.getPath()
	_, time := p.findPortInBytes(path, port)
	return time
}

func (p *Port) findPortInBytes(path string, port string) (int, time.Time) {
	indexRow := -1
	var oldestTime time.Time
	lines := p.getAllLines(path)

	for i, row := range lines {
		fields := strings.Fields(row)
		t, portExist := fields[0], fields[1]
		if portExist == port {
			indexRow = i
			unixInt, _ := strconv.ParseInt(t, 10, 64)
			oldestTime = time.Unix(unixInt, 0)
			break
		}
	}
	return indexRow, oldestTime
}

func (p *Port) store(port string) {
	path := p.getPath()
	idx, _ := p.findPortInBytes(path, port)
	lines := p.getAllLines(path)

	line := fmt.Sprintf("%d %s", time.Now().Unix(), port)

	if idx != -1 {
		lines[idx] = line
	} else {
		lines = append(lines, line)
	}
	output := strings.Join(lines, "\n")
	os.WriteFile(path, []byte(output), 0o644)
}

func (p *Port) getAllLines(path string) []string {
	input, _ := os.ReadFile(path)
	lines := make([]string, 0)
	parts := strings.Split(string(input), "\n")
	for _, part := range parts {
		if part != "" {
			lines = append(lines, part)
		}
	}
	return lines
}

func (p *Port) getPath() string {
	root, _ := utils.GetRootOfProject()

	return filepath.Join(root, "ports")
}
