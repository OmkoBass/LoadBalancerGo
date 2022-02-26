package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	// Need a mutex when updating the
	// lastForwardIndex variable
	mutex              sync.Mutex
	lastForwardedIndex = 0

	totalBackends = 4

	serverListFirst = []*server{
		newServer("Server-1", "http://localhost:8000"),
		newServer("Server-2", "http://localhost:8001"),
		newServer("Server-3", "http://localhost:8002"),
		newServer("Server-4", "http://localhost:8003"),
	}

	serverListSecond = []*server{
		newServer("Server-4", "http://localhost:8004"),
		newServer("Server-5", "http://localhost:8005"),
		newServer("Server-6", "http://localhost:8006"),
		newServer("Server-7", "http://localhost:8007"),
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
	//fmt.Println(server.URL)
	req.Host = server.URL
	server.ReverseProxy.ServeHTTP(res, req)
}

func getBackend() *server {
	// This is a problem for my usecase
	// I need to use different servers
	// dependant on the current date
	// First half of the month i will be
	// using serverListFirst
	// second half of the month i will be
	// using serverListSecond

	currentTime := time.Now()

	if currentTime.Day() < 15 {
		backend := serverListFirst[lastForwardedIndex]

		mutex.Lock()
		lastForwardedIndex = (lastForwardedIndex + 1) % totalBackends
		mutex.Unlock()

		fmt.Println(backend.Name)
		return backend
	}

	backend := serverListSecond[lastForwardedIndex]

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
