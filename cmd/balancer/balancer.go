package balancer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/sebasfalcone/go-load-balancer/internal/pkg/endpoint"
)

// Config is the configuration of the load balancer (proxy and backends)
type Config struct {
	Proxy    Proxy              `json:"proxy"`    // Proxy configuration
	Backends []endpoint.Backend `json:"backends"` // List of backends
}

// Proxy configuration
type Proxy struct {
	Port string `json:"port"` // Port to listen on
}

func healthCheck() {
	var statusMessage = map[bool]string{false: "dead", true: "alive"}
	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-ticker.C:
			for _, backend := range config.Backends {
				status := backend.UpdateAliveStatus()
				log.Printf("URL: %v - Status: %v", backend.URL, statusMessage[status])
			}
		}
	}

}

var config Config
var backendIdx int
var backendMaxIdx int
var mutex sync.Mutex

// Start the load balancer
func Start() {

	loadConfigurations()

	// Execute health check concurrently
	go healthCheck()

	server := http.Server{
		Addr:    ":" + config.Proxy.Port,
		Handler: http.HandlerFunc(loadBalancerHandler),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

func loadConfigurations() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	json.Unmarshal(data, &config)
	backendMaxIdx = len(config.Backends)

	for index := 0; index < backendMaxIdx; index++ {
		// Validate configurations and initialize alive status
		backend := config.Backends[index]
		backend.UpdateAliveStatus()
	}
}

func loadBalancerHandler(writter http.ResponseWriter, request *http.Request) {
	mutex.Lock()
	targetURL := getTargetUrl()
	mutex.Unlock()

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Println("Unable to redirect trafic to ", targetURL, ". Error:", err)
		loadBalancerHandler(w, r)
	}

	reverseProxy.ServeHTTP(writter, request)
}

// Returns the next endpoint URL to redirect the request
func getTargetUrl() *url.URL {

	// Round Robin
	backend := config.Backends[backendIdx]
	for count := 0; count < backendMaxIdx; count++ {

		backendIdx++
		if backendIdx >= backendMaxIdx {
			backendIdx = 0
		}

		if backend.UpdateAliveStatus() != false {
			break
		}

		backend = config.Backends[backendIdx]
	}

	targetURL, err := url.Parse(backend.URL)
	if err != nil {
		log.Fatal(err.Error())
	}

	return targetURL
}
