package worker

import "momo/entity"

type MockWorkerProxy struct{}

func (mwp *MockWorkerProxy) GetAvailablePorts(requiredPorts uint32, portsUsed []string) ([]string, error) {
	ans := []string{}
	for i := 0; i < int(requiredPorts); i++ {
		ans = append(ans, "1234")
	}
	return ans, nil
}

func (mwp *MockWorkerProxy) GetMetric() (uint32, string, error) {
	return 10, string(entity.High), nil
}

func (mwp *MockWorkerProxy) Close() {}
