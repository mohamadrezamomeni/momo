package adapter

import (
	"momo/entity"
	workerProxy "momo/proxy/worker"
)

type WorkerProxy interface {
	Close()
	GetAvailablePorts(uint32, []string) ([]string, error)
	GetMetric() (uint32, entity.HostStatus, error)
}

func AdaptedWorkerFactory(address string, port string) (WorkerProxy, error) {
	wp, err := workerProxy.New(&workerProxy.Config{
		Address: address,
		Port:    port,
	})
	if err != nil {
		return nil, err
	}
	return wp, nil
}
