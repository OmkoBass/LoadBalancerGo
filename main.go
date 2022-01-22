package main

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	// Need a mutex when updating the
	// lastForwardIndex variable
	mutex              sync.Mutex
	lastForwardedIndex = 0

	totalBackends = 5
	backendList   = []*server{
		newServer("Flask-1", "http://localhost:5000"),
		newServer("Flask-2", "http://localhost:5001"),
		newServer("Flask-3", "http://localhost:5002"),
		newServer("Flask-4", "http://localhost:5003"),
		newServer("Flask-5", "http://localhost:5004"),
	}
)

func main() {
	http.HandleFunc("/", handleRequest)

	go startHealthCheck(5)

	http.ListenAndServe(":8000", nil)
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	// Find the alive server
	server, err := getAliveServer()

	// If there's not alive servers
	// send a response
	if err != nil {
		http.Error(res, "Request cannot be handled. Reason: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	server.ReverseProxy.ServeHTTP(res, req)
}

func getBackend() *server {
	backend := backendList[lastForwardedIndex]

	mutex.Lock()
	lastForwardedIndex = (lastForwardedIndex + 1) % totalBackends
	mutex.Unlock()

	return backend
}

func getAliveServer() (*server, error) {
	// n servers <==> loop n times
	for i := 0; i < totalBackends; i++ {
		server := getBackend()

		if server.Alive {
			return server, nil
		}
	}

	return nil, fmt.Errorf("all servers are dead")
}
