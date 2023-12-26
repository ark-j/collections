package reqwest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
)

type Reqwest struct {
	client *http.Client           // http client
	trace  bool                   // to enable tracing or not
	tracer *httptrace.ClientTrace // actual tracer implmentation
}

func New(headers map[string]string, trace bool) *Reqwest {
	return (&Reqwest{
		client: &http.Client{},
		trace:  trace,
		tracer: getTracer(),
	}).SetTransport(defaultTransport)
}

// SetTransport set the httptransport,
// if povided transport is nil,
// default transport will be used.
func (r *Reqwest) SetTransport(t *http.Transport) *Reqwest {
	if t != nil {
		r.client.Transport = t
	}
	return r
}

// DisableRedirect disable the redirects in http.Client.
// By default redirect are not disabled and
// follows upto configured redirects in http client.
func (r *Reqwest) DisableRedirect() *Reqwest {
	r.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return r
}

// SetCookieJar set cookie jar with contained cookies
// by default no cookie jar is setup
func (r *Reqwest) SetCookieJar(jar http.CookieJar) *Reqwest {
	r.client.Jar = jar
	return r
}

// SetTracer replace default tracer with your own implementation.
func (r *Reqwest) SetTracer(tracer *httptrace.ClientTrace) *Reqwest {
	if tracer != nil {
		r.tracer = tracer
	}
	return r
}

// Get is http get method
func (r *Reqwest) Get(ctx context.Context, uri string, headers map[string]string) (*http.Response, error) {
	return r.request(ctx, http.MethodGet, uri, nil, headers)
}

// Head is http head method follows upto 10 redirect
func (r *Reqwest) Head(ctx context.Context, uri string, headers map[string]string) (*http.Response, error) {
	return r.request(ctx, http.MethodHead, uri, nil, headers)
}

// Post is http post method
func (r *Reqwest) Post(ctx context.Context, uri string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return r.request(ctx, http.MethodPost, uri, body, headers)
}

// Put is http put method
func (r *Reqwest) Put(ctx context.Context, uri string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return r.request(ctx, http.MethodPut, uri, body, headers)
}

// Patch is http patch method
func (r *Reqwest) Patch(ctx context.Context, uri string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return r.request(ctx, http.MethodPatch, uri, body, headers)
}

// Delete is http delete method
func (r *Reqwest) Delete(ctx context.Context, uri string, headers map[string]string) (*http.Response, error) {
	return r.request(ctx, http.MethodDelete, uri, nil, headers)
}

// request is lowlevel function to perform request in Get, Put, Post, Head, Patch, Delete
func (r *Reqwest) request(ctx context.Context, method, uri string, body io.Reader, headers map[string]string) (*http.Response, error) {
	if r.trace {
		ctx = httptrace.WithClientTrace(ctx, r.tracer)
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	return r.client.Do(req)
}
