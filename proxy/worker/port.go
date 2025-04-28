package worker

import (
	"context"
	"time"

	pb "momo/contract/gogrpc/port"
	momoError "momo/pkg/error"
)

func (pw *ProxyWorker) GetAvailablePorts(requestNumberOfPorts uint32, portsUsed []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := pw.portClient.GetAvailablePorts(ctx, &pb.PortAssignRequest{
		RequestNumberOfPort: requestNumberOfPorts,
		PortsUsed:           portsUsed,
	})
	if err != nil {
		return nil, momoError.Errorf("faild to request the error was %v", err)
	}
	return res.Ports, nil
}
