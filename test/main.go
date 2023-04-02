package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/sebasfalcone/go-load-balancer/cmd/balancer"
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
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(3)

	go func() {
		serveBackend("web1", ":81")
		waitGroup.Done()
	}()

	go func() {
		serveBackend("web2", ":82")
		waitGroup.Done()
	}()

	go func() {
		serveBackend("web3", ":83")
		waitGroup.Done()
	}()

	balancer.Start()
	waitGroup.Wait()
}
