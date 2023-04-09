package mocks

import (
	"net/url"
)

type HealthMock struct {
	// The origin of the proxy
	origin *url.URL

	// The health check function and period
	check func(addr *url.URL) bool

	// The health status of the proxy
	alive bool
}

func (h *HealthMock) CheckFunction(addr *url.URL) bool {
	if addr.Host == "valid" {
		return true
	} else {
		return false
	}
}

func NewHealth(origin *url.URL) *HealthMock {
	h := &HealthMock{
		origin: origin,
	}

	h.check = h.CheckFunction
	h.alive = h.check(h.origin)

	return h
}
