package worker

import (
	"github.com/mohamadrezamomeni/momo/entity"
)

type MockWorkerProxy struct{}

func (mwp *MockWorkerProxy) GetAvailablePorts(requiredPorts uint32, portsUsed []string) ([]string, error) {
	ans := []string{}
	for i := 0; i < int(requiredPorts); i++ {
		ans = append(ans, "1234")
	}
	return ans, nil
}

func (mwp *MockWorkerProxy) OpenPorts(ports []string) ([]string, error) {
	portsFailed := make([]string, 0)
	for _, port := range ports {
		if port == "3333" || port == "5555" || port == "8888" {
			portsFailed = append(portsFailed, port)
		}
	}
	return portsFailed, nil
}

func (mwp *MockWorkerProxy) GetMetric() (uint32, entity.HostStatus, error) {
	return 10, entity.High, nil
}

func (mwp *MockWorkerProxy) Close() {}
