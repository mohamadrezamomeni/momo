package worker

import (
	"context"
	"time"

	pb "github.com/mohamadrezamomeni/momo/contract/gogrpc/port"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func (pw *ProxyWorker) GetAvailablePorts(requestNumberOfPorts uint32, portsUsed []string) ([]string, error) {
	scope := "workerProxy.GetAvailablePorts"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := pw.portClient.GetAvailablePorts(ctx, &pb.PortAssignRequest{
		RequestNumberOfPort: requestNumberOfPorts,
		PortsUsed:           portsUsed,
	})
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(requestNumberOfPorts, portsUsed).Errorf(
			"the requiestNumberOfPorts is %d and portsUsed is %+v and address of worker is %s",
			requestNumberOfPorts,
			portsUsed,
			pw.address,
		)
	}
	return res.Ports, nil
}

func (pw *ProxyWorker) OpenPorts(ports []string) ([]string, error) {
	scope := "workerProxy.OpenPort"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := pw.portClient.OpenPorts(ctx, &pb.OpenPortsRequest{
		Ports: ports,
	})
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).Input(ports).ErrorWrite()
	}
	return res.Ports, nil
}
