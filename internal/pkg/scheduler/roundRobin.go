package scheduler

import (
	"fmt"
	"sync/atomic"

	"github.com/sebasfalcone/go-load-balancer/internal/pkg/proxy"
)

type RoundRobin struct {
	// The list of proxys to be used
	proxys []*proxy.Proxy

	// The current index of the proxy to be used
	current int32

	// The amount of proxys
	amount int32
}

func (r *RoundRobin) AddProxys(p ...*proxy.Proxy) Scheduler {
	r.proxys = append(r.proxys, p...)
	r.amount = int32(len(r.proxys))
	r.current = -1

	return r
}

func (r *RoundRobin) Instance() Scheduler {
	return r
}

func (r *RoundRobin) Next() (*proxy.Proxy, error) {

	for count := int32(0); count < r.amount; count++ {

		idx := atomic.AddInt32(&r.current, 1) % r.amount
		atomic.StoreInt32(&r.current, idx)

		if r.proxys[idx].Status() {
			return r.proxys[idx], nil
		}
	}

	return nil, fmt.Errorf("all proxys are unavailable")
}
