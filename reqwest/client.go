package reqwest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
)

// Change this to desired user agent header
// or you can always overwirte using Header Options
const HeaderUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type Reqwest struct {
	client *http.Client
	tracer *httptrace.ClientTrace
	trace  bool
}

func New(trace bool) *Reqwest {
	return (&Reqwest{
		client: &http.Client{},
		trace:  trace,
		tracer: getTracer(),
	}).SetTransport(defaultTransport)
}

// SetTransport set the httptransport,
// if povided transport is nil,
// default transport will be used.
func (r *Reqwest) SetTransport(t http.RoundTripper) *Reqwest {
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
func (r *Reqwest) Get(ctx context.Context, uri string, opts ...Options) (*http.Response, error) {
	return request(ctx, r, http.MethodGet, uri, nil, opts...)
}

// Head is http head method follows upto 10 redirect
func (r *Reqwest) Head(ctx context.Context, uri string, opts ...Options) (*http.Response, error) {
	return request(ctx, r, http.MethodHead, uri, nil, opts...)
}

// Post is http post method
func (r *Reqwest) Post(ctx context.Context, uri string, body io.Reader, opts ...Options) (*http.Response, error) {
	return request(ctx, r, http.MethodPost, uri, body, opts...)
}

// Put is http put method
func (r *Reqwest) Put(ctx context.Context, uri string, body io.Reader, opts ...Options) (*http.Response, error) {
	return request(ctx, r, http.MethodPut, uri, body, opts...)
}

// Patch is http patch method
func (r *Reqwest) Patch(ctx context.Context, uri string, body io.Reader, opts ...Options) (*http.Response, error) {
	return request(ctx, r, http.MethodPatch, uri, body, opts...)
}

// Delete is http delete method
func (r *Reqwest) Delete(ctx context.Context, uri string, opts ...Options) (*http.Response, error) {
	return request(ctx, r, http.MethodDelete, uri, nil, opts...)
}

// request is lowlevel function to perform request in Get, Put, Post, Head, Patch, Delete
func request(
	ctx context.Context,
	r *Reqwest,
	method,
	uri string,
	body io.Reader,
	opts ...Options,
) (*http.Response, error) {
	// if trace is available
	if r.trace {
		ctx = httptrace.WithClientTrace(ctx, r.tracer)
	}

	// initiate options for headers and queries
	ho := &HTTPOptions{}
	for _, o := range opts {
		o(ho)
	}

	// initiate request with context
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	// initiate request header for general uses
	req.Header.Set("User-Agent", HeaderUserAgent)

	// set all optional headers
	if len(ho.headers) > 0 {
		for k, v := range ho.headers {
			req.Header.Set(k, v)
		}
	}
	// set all optional queries
	if len(ho.queries) > 0 {
		q := req.URL.Query()
		for k, v := range ho.queries {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return r.client.Do(req)
}
