package vpn

import (
	"github.com/mohamadrezamomeni/momo/adapter"
	mockProxyVPN "github.com/mohamadrezamomeni/momo/mocks/proxy/vpn"
)

type MockVPN struct{}

func New() *MockVPN {
	return &MockVPN{}
}

func (mv *MockVPN) MakeProxy() (adapter.ProxyVPN, error) {
	return &mockProxyVPN.MockProxy{}, nil
}
