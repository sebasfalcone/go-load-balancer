package health

import (
	"net"
	"net/url"
	"sync"
	"time"
)

type Health struct {
	// The origin of the proxy
	origin *url.URL

	// The mutex to protect the health check
	mu sync.RWMutex

	// The health check function and period
	check func(addr *url.URL) bool

	// The health check period
	period time.Duration

	// The health check cancel channel
	cancel chan struct{}

	// The health status of the proxy
	alive bool
}

var defaultCheckTimeout = 10 * time.Second
var defaultCheckPeriod = 10 * time.Second

// Creates a new health check for the proxy
func NewHealth(origin *url.URL) *Health {

	h := &Health{
		origin: origin,
		check:  defaultCheckFunction,
		period: defaultCheckPeriod,
		cancel: make(chan struct{}),
		alive:  defaultCheckFunction(origin),
	}

	h.start()
	return h
}

// Default health check function
var defaultCheckFunction = func(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, defaultCheckTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// Returns the health status of the proxy
func (h *Health) Status() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.alive
}

// Sets the health check function and period, for the proxy health check
func (h *Health) SetHealthCheck(checkFunction func(addr *url.URL) bool, period time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.stop()
	h.check = checkFunction
	h.period = period
	h.cancel = make(chan struct{})
	h.alive = h.check(h.origin)
	h.start()
}

// Starts the health check
func (h *Health) start() {

	// Check the health of the proxy
	checkHealth := func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		alive := h.check(h.origin)
		h.alive = alive
	}

	// Start the health check ticker
	go func() {
		ticker := time.NewTicker(h.period)
		for {
			select {
			case <-ticker.C:
				checkHealth()
			case <-h.cancel:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the health check
func (h *Health) stop() {
	if h.cancel != nil {
		h.cancel <- struct{}{}
		close(h.cancel)
		h.cancel = nil
	}
}
