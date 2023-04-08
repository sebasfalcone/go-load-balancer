package loadBalancer

import (
	"log"
	"net/http"

	"github.com/sebasfalcone/go-load-balancer/internal/pkg/scheduler"
)

type LoadBalancer struct {
	// The scheduler to be used
	sched scheduler.Scheduler
}

// Initialize starts the load balancer
func Initialize(s scheduler.Scheduler, addr string) {

	log.Printf("load balancer started at port %s", addr)

	lb := LoadBalancer{sched: s}

	handler := func(w http.ResponseWriter, r *http.Request) {
		proxy, err := lb.sched.Next()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Redirecting request to %s", proxy.GetName())
		proxy.ServeHTTP(w, r)
	}

	server := http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(handler),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

}
