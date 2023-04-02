package endpoint

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

var mutex sync.RWMutex

// Backend configuration and status
type Backend struct {
	URL   string `json:"url"` // URL of the backend
	alive bool   // true if backend is alive
}

// Checks the backend status, updates it and returns the current status (true if alive)
func (backend *Backend) UpdateAliveStatus() bool {
	mutex.Lock()
	defer mutex.Unlock()

	targetURL, err := url.Parse(backend.URL)
	if err != nil {
		log.Fatal(err.Error())
	}

	connection, err := net.DialTimeout("tcp", targetURL.Host, time.Minute*1)
	if err != nil {
		log.Printf("%v is unreachable. Error: %v", targetURL, err.Error())
		backend.alive = false
	} else {
		connection.Close()
		backend.alive = true
	}
	fmt.Println("Backend: ", targetURL, ". Status: ", backend.alive)
	return backend.alive
}

// Returns the current status of the backend (true if alive)
func (backend *Backend) AliveStatus() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	target, _ := url.Parse(backend.URL)
	fmt.Println("Backend: ", target, ". Status: ", backend.alive)
	return backend.alive
}
