package worker

import (
	"context"
	"time"

	pb "momo/contract/gogrpc/port"
	momoError "momo/pkg/error"
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
		return nil, momoError.Wrap(err).Scope(scope).Errorf(
			"the requiestNumberOfPorts is %d and portsUsed is %+v and address of worker is %s",
			requestNumberOfPorts,
			portsUsed,
			pw.address,
		)
	}
	return res.Ports, nil
}
