package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sebasfalcone/go-load-balancer/internal/pkg/loadBalancer"
	"github.com/sebasfalcone/go-load-balancer/internal/pkg/proxy"
	"github.com/sebasfalcone/go-load-balancer/internal/pkg/scheduler"
)

func serveBackend(serverName string, serverPort string) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(writter http.ResponseWriter, request *http.Request) {
		writter.WriteHeader(http.StatusOK)
		fmt.Fprintln(writter, "Server: ", serverName)
		fmt.Fprintln(writter, "Response: ", request.Header)
	}))

	http.ListenAndServe(serverPort, mux)
}

func main() {

	go func() {
		serveBackend("web1", ":81")
	}()

	go func() {
		serveBackend("web2", ":82")
	}()

	go func() {
		serveBackend("web3", ":83")
	}()

	p1Url, _ := url.Parse("http://localhost:81")
	p1 := proxy.NewProxy(p1Url, "web1")

	p2Url, _ := url.Parse("http://localhost:82")
	p2 := proxy.NewProxy(p2Url, "web2")

	p3Url, _ := url.Parse("http://localhost:83")
	p3 := proxy.NewProxy(p3Url, "web3")

	r := scheduler.RoundRobin{}
	r.AddProxys(p1, p2, p3)

	loadBalancer.Initialize(r.Instance(), ":80")
}
