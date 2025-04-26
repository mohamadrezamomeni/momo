package worker

import (
	"fmt"

	"google.golang.org/grpc"
)

type ProxyWorker struct {
	conn    *grpc.ClientConn
	address string
}

func New(cfg *Config) (*ProxyWorker, error) {
	address := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return &ProxyWorker{}, err
	}

	return &ProxyWorker{
		conn:    conn,
		address: address,
	}, nil
}

func (ps *ProxyWorker) Close() {
	ps.conn.Close()
}
