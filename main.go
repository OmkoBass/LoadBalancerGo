package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	mutex              sync.Mutex
	lastForwardedIndex = 0

	totalBackends = 5
	backendList   = []*httputil.ReverseProxy{
		createHost("http://localhost:5000/alive"),
		createHost("http://localhost:5001/alive"),
		createHost("http://localhost:5002/alive"),
		createHost("http://localhost:5003/alive"),
		createHost("http://localhost:5004/alive"),
	}
)

func main() {
	http.HandleFunc("/", handleRequest)

	http.ListenAndServe(":8000", nil)
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	server := getBackend()

	server.ServeHTTP(res, req)
}

func getBackend() *httputil.ReverseProxy {
	backend := backendList[lastForwardedIndex]

	mutex.Lock()
	lastForwardedIndex = (lastForwardedIndex + 1) % totalBackends
	mutex.Unlock()

	return backend
}

func createHost(urlToParse string) *httputil.ReverseProxy {
	parsedUrl, _ := url.Parse(urlToParse)
	return httputil.NewSingleHostReverseProxy(parsedUrl)
}
