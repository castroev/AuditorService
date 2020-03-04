package customclient

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// ClientSettings is a struct to hold the timeout values we want to be able to set on our HTTP client
type ClientSettings struct {
	Connect          time.Duration
	ConnKeepAlive    time.Duration
	ExpectContinue   time.Duration
	IdleConn         time.Duration
	MaxAllIdleConns  int
	MaxHostIdleConns int
	ResponseHeader   time.Duration
	TLSHandshake     time.Duration
}

// NewHTTPClientWithTimeouts creates a custom http client with specified transport settings
func NewHTTPClientWithTimeouts(httpSettings ClientSettings) *http.Client {
	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeader,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: httpSettings.ConnKeepAlive,
			DualStack: true,
			Timeout:   httpSettings.Connect,
		}).DialContext,
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConn,
		TLSHandshakeTimeout:   httpSettings.TLSHandshake,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinue,
	}

	// So client makes HTTP/2 requests
	http2.ConfigureTransport(tr)

	return &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}
}
