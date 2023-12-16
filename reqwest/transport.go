package reqwest

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

var defaultTransport = &http.Transport{
	DialContext: transportDailContext(),
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	MaxIdleConns:          maxIdleConns,
	MaxIdleConnsPerHost:   maxIdleConnsPerHost,
	IdleConnTimeout:       idleConnTimeout,
	ExpectContinueTimeout: expectContinueTimeout,
	TLSHandshakeTimeout:   tlsHandshakeTimeout,
	ForceAttemptHTTP2:     true,
	WriteBufferSize:       mb,
	ReadBufferSize:        mb,
}

// SetProxy set proxy to defaultTransport
// if you're using custom transport it is assumed that you have provide proxy with it
func SetProxy(proxy func(r *http.Request) (*url.URL, error)) {
	defaultTransport.Proxy = proxy
}

// GetDefaultTransport returns [*http.Transport] which you can configure to your like other than defaults
func GetDefaultTransport() *http.Transport {
	return defaultTransport
}

// transportDailContext return DailContext Func for setting it in transport
// usable for field such as DialContext and DialTLSContext
func transportDailContext() func(context.Context, string, string) (net.Conn, error) {
	return (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
}
