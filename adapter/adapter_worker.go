package adapter

import (
	"github.com/mohamadrezamomeni/momo/entity"
	workerProxy "github.com/mohamadrezamomeni/momo/proxy/worker"
)

type WorkerProxy interface {
	Close()
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
