package client

import (
	"net"
	"net/http"
	"time"
)

const (
	DialTimeout           = 5 * time.Second
	DialKeepAlive         = 30 * time.Second
	TLSHandshakeTimeout   = 5 * time.Second
	ResponseHeaderTimeout = 30 * time.Second
	MaxIdleConnsPerHost   = 100
)

// NewHTTPClient creates the shared transport used by provider clients.
func NewHTTPClient() *http.Client {
	return &http.Client{Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: DialTimeout, KeepAlive: DialKeepAlive}).DialContext,
		TLSHandshakeTimeout:   TLSHandshakeTimeout,
		ResponseHeaderTimeout: ResponseHeaderTimeout,
		MaxIdleConnsPerHost:   MaxIdleConnsPerHost,
	}}
}
