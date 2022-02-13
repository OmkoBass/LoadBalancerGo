package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type server struct {
	Name         string
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Alive        bool
}

func newServer(name string, urlToParse string) *server {
	parsedUrl, _ := url.Parse(urlToParse)
	fmt.Println(urlToParse)

	targetQuery := parsedUrl.RawQuery
	reverseProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
            req.Host = parsedUrl.Host
			req.URL.Scheme = parsedUrl.Scheme
            req.URL.Host = parsedUrl.Host
            req.URL.Path = singleJoiningSlash(parsedUrl.Path, req.URL.Path)
			if targetQuery == "" || req.URL.RawQuery == "" {
                req.URL.RawQuery = targetQuery + req.URL.RawQuery
            } else {
                req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
            }
        },
	}

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

func singleJoiningSlash(a, b string) string {
    aslash := strings.HasSuffix(a, "/")
    bslash := strings.HasPrefix(b, "/")
    switch {
    case aslash && bslash:
        return a + b[1:]
    case !aslash && !bslash:
        return a + "/" + b
    }
    return a + b
}
