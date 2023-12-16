package reqwest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"
)

const (
	mb                    = 1024 * 1024
	maxIdleConnsPerHost   = 2
	idleConnTimeout       = 2 * time.Minute
	expectContinueTimeout = 1 * time.Second
	tlsHandshakeTimeout   = 10 * time.Second
	maxIdleConns          = 512
)

// TODO: add default headers and option to add headers in request
type Reqwest struct {
	client  *http.Client
	trace   bool
	tracer  *httptrace.ClientTrace
	headers map[string]string
	mu      *sync.RWMutex
}

func New(headers map[string]string, trace bool) *Reqwest {
	return (&Reqwest{
		&http.Client{},
		trace,
		getTracer(),
		headers,
		&sync.RWMutex{},
	}).SetTransport(nil)
}

// SetTransport set the httptransport if povided transport is nil,
// it will set default transport
func (r *Reqwest) SetTransport(t *http.Transport) *Reqwest {
	if t != nil {
		r.client.Transport = t
	} else {
		r.client.Transport = defaultTransport
	}
	return r
}

// SetCookieJar set cookie jar with contained cookies
// by default no cookie jar is setup
func (r *Reqwest) SetCookieJar(jar http.CookieJar) *Reqwest {
	r.client.Jar = jar
	return r
}

// GetClient return pointer to underlying [net/http.Client].
// you can use this client for low level operations
func (r *Reqwest) GetClient() *http.Client {
	return r.client
}

// Get is http get method
func (r *Reqwest) Get(ctx context.Context, uri string) (*http.Response, error) {
	return r.request(ctx, http.MethodGet, uri, nil)
}

// Head is http head method follows upto 10 redirect
func (r *Reqwest) Head(ctx context.Context, uri string) (*http.Response, error) {
	return r.request(ctx, http.MethodHead, uri, nil)
}

// Post is http post method
func (r *Reqwest) Post(ctx context.Context, uri string, body io.Reader) (*http.Response, error) {
	return r.request(ctx, http.MethodPost, uri, body)
}

// Put is http put method
func (r *Reqwest) Put(ctx context.Context, uri string, body io.Reader) (*http.Response, error) {
	return r.request(ctx, http.MethodPut, uri, body)
}

// request is lowlevel function to perform request in Get, Put, Post, Head
func (r *Reqwest) request(ctx context.Context, method, uri string, body io.Reader) (*http.Response, error) {
	if r.trace {
		ctx = httptrace.WithClientTrace(ctx, r.tracer)
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	r.mu.RLock()
	if len(r.headers) > 0 {
		for k, v := range r.headers {
			req.Header.Add(k, v)
		}
	}
	r.mu.RUnlock()
	return r.client.Do(req)
}
