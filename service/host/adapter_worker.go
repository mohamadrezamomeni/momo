package host

import workerProxy "momo/proxy/worker"

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
