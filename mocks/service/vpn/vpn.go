package vpn

import (
	proxyvpn "momo/proxy/vpn"

	mockProxyVPN "momo/mocks/proxy/vpn"
)

type MockVPN struct{}

func New() *MockVPN {
	return &MockVPN{}
}

func (mv *MockVPN) MakeProxy() (proxyvpn.IProxyVPN, error) {
	return &mockProxyVPN.MockProxy{}, nil
}
