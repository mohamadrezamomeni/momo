package vpn

import (
	"momo/adapter"
	mockProxyVPN "momo/mocks/proxy/vpn"
)

type MockVPN struct{}

func New() *MockVPN {
	return &MockVPN{}
}

func (mv *MockVPN) MakeProxy() (adapter.ProxyVPN, error) {
	return &mockProxyVPN.MockProxy{}, nil
}
