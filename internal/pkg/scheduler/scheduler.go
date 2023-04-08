package scheduler

import proxy "github.com/sebasfalcone/go-load-balancer/internal/pkg/proxy"

type Scheduler interface {
	// Returns the next proxy to be used or an error if all the endpoints where unavailable
	Next() (*proxy.Proxy, error)

	// Adds a new proxy/s to the scheduler
	AddProxys(proxy ...*proxy.Proxy) Scheduler

	// Gets the scheduler instance
	Instance() Scheduler
}
