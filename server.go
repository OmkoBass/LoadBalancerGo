package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type server struct {
	Name         string
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Alive        bool
}

func newServer(name string, urlToParse string) *server {
	parsedUrl, _ := url.Parse(urlToParse)

	reverseProxy := httputil.NewSingleHostReverseProxy(parsedUrl)

	return &server{
		Name:         name,
		URL:          urlToParse,
		ReverseProxy: reverseProxy,
		Alive:        true,
	}
}

func (server *server) isAlive() bool {
	// Checks if the server responds
	// if it does with a 200 then it's alive
	// if it doesn't respond or it's not 200
	// then it's dead

	resp, err := http.Head(server.URL + "/alive")

	server.Alive = true

	if err != nil {
		server.Alive = false
		return server.Alive
	}

	if resp.StatusCode != http.StatusOK {
		server.Alive = false
		return server.Alive
	}

	return server.Alive
}
