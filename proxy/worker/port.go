package worker

import momoError "momo/pkg/error"

func (pw *ProxyWorker) GetAvailablePort() (string, error) {
	port, err := pw.GetAvailablePort()
	if err != nil {
		momoError.Errorf("error has happend the port hasn't assigned the problem was %v", err)
	}
	return port, nil
}
