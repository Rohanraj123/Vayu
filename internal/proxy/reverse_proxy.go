package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	backendURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(backendURL), nil
}

func ProxyHandler(target string) (http.HandlerFunc, error) {
	proxy, err := NewReverseProxy(target)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}, nil
}
