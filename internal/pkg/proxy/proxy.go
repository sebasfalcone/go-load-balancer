package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/sebasfalcone/go-load-balancer/internal/pkg/proxy/health"
)

// Proxy configuration and status
type Proxy struct {
	health *health.Health
	proxy  *httputil.ReverseProxy
	load   int32
	name   string
}

// Creates a new proxy
func NewProxy(addr *url.URL, name string) *Proxy {
	log.Printf("creating proxy for %s", addr)
	return &Proxy{
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		health: health.NewHealth(addr),
		load:   0,
		name:   name,
	}
}

// Serves the proxy
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&p.load, 1)
	defer atomic.AddInt32(&p.load, -1)

	log.Printf("Request body: %s", r.Body)
	p.proxy.ServeHTTP(w, r)
}

// Returns the proxy status
// True if the proxy is available, false otherwise
func (p *Proxy) Status() bool {
	return p.health.Status()
}

// Sets the health check function and period, for the proxy health check
func (p *Proxy) SetHealthCheck(checkFunction func(addr *url.URL) bool, period time.Duration) {
	p.health.SetHealthCheck(checkFunction, period)
}

// Returns the proxy load
func (p *Proxy) GetLoad() int32 {
	return atomic.LoadInt32(&p.load)
}

// Returns the proxy name
func (p *Proxy) GetName() string {
	return p.name
}
