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

	totalBackends = 4
	backendList   = []*server{
		newServer("Story-Saver", "https://private-story-saver.herokuapp.com"),
		newServer("Story-Saver-1", "https://private-story-saver-1.herokuapp.com"),
		newServer("Story-Saver-2", "https://private-story-saver-2.herokuapp.com"),
		newServer("Story-Saver-3", "https://private-story-saver-3.herokuapp.com"),
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
