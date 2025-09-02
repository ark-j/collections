package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
)

// Change this to desired user agent header
// or you can always overwrite using Header Options
const HeaderUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"

type Client struct {
	client *http.Client
	tracer *httptrace.ClientTrace
	trace  bool
}

func New(trace bool) *Client {
	return (&Client{
		client: &http.Client{},
		trace:  trace,
		tracer: getTracer(),
	}).SetTransport(defaultTransport)
}

// SetTransport set the httptransport,
// if povided transport is nil,
// default transport will be used.
func (c *Client) SetTransport(t http.RoundTripper) *Client {
	if t != nil {
		c.client.Transport = t
	}
	return c
}

// DisableRedirect disable the redirects in http.Client.
// By default redirect are not disabled and
// follows upto configured redirects in http client.
func (c *Client) DisableRedirect() *Client {
	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return c
}

// SetCookieJar set cookie jar with contained cookies
// by default no cookie jar is setup
func (c *Client) SetCookieJar(jar http.CookieJar) *Client {
	c.client.Jar = jar
	return c
}

// SetTracer replace default tracer with your own implementation.
func (c *Client) SetTracer(tracer *httptrace.ClientTrace) *Client {
	if tracer != nil {
		c.tracer = tracer
	}
	return c
}

// Get is http get method
func (c *Client) Get(ctx context.Context, uri string, ho *HTTPOptions) (*http.Response, error) {
	return c.Exec(ctx, http.MethodGet, uri, nil, ho)
}

// Head is http head method follows upto 10 redirect
func (c *Client) Head(ctx context.Context, uri string, ho *HTTPOptions) (*http.Response, error) {
	return c.Exec(ctx, http.MethodHead, uri, nil, ho)
}

// Post is http post method
func (c *Client) Post(
	ctx context.Context,
	uri string,
	body io.Reader,
	ho *HTTPOptions,
) (*http.Response, error) {
	return c.Exec(ctx, http.MethodPost, uri, body, ho)
}

// Put is http put method
func (c *Client) Put(
	ctx context.Context,
	uri string,
	body io.Reader,
	ho *HTTPOptions,
) (*http.Response, error) {
	return c.Exec(ctx, http.MethodPut, uri, body, ho)
}

// Patch is http patch method
func (c *Client) Patch(
	ctx context.Context,
	uri string,
	body io.Reader,
	ho *HTTPOptions,
) (*http.Response, error) {
	return c.Exec(ctx, http.MethodPatch, uri, body, ho)
}

// Delete is http delete method
func (c *Client) Delete(ctx context.Context, uri string, ho *HTTPOptions) (*http.Response, error) {
	return c.Exec(ctx, http.MethodDelete, uri, nil, ho)
}

// Exec performs the HTTP request with the given method, uri, and options.
//
// Hook execution order:
//
//  1. requestHook — runs before sending the request.
//  2. retryHook   — if defined, takes full control over retries and
//     determines the final response. In this case, responseHook is
//     NOT invoked.
//  3. responseHook — runs only if no retryHook is defined.
//
// Important:
//
//   - If retryHook is defined (custom or default), responseHook will be ignored.
//     This avoids conflicts from reading res.Body multiple times.
//
//   - When using the default retryHook, place any post-processing logic
//     (e.g. decoding JSON, logging, validation) in the Cond function itself.
//
//   - When writing a custom retryHook, encapsulate your retry decision and
//     any post-processing logic inside the retryHook implementation.
//
// This ensures hooks remain predictable and prevents accidental multiple
// reads of the response body.
func (c *Client) Exec(
	ctx context.Context,
	method, uri string,
	body io.Reader,
	ho *HTTPOptions,
) (*http.Response, error) {
	// if trace is available
	if c.trace {
		ctx = httptrace.WithClientTrace(ctx, c.tracer)
	}

	if ho == nil {
		ho = &HTTPOptions{}
	}

	// initiate request with context
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	// initiate request header for general uses
	req.Header.Set("User-Agent", HeaderUserAgent)

	// set all optional headers
	for k, v := range ho.headers {
		req.Header.Set(k, v)
	}

	// set all optional queries
	q := req.URL.Query()
	for k, v := range ho.queries {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	if ho.requestHook != nil {
		if err := ho.requestHook(req); err != nil {
			return nil, fmt.Errorf("failed to execute request hook: %w", err)
		}
	}

	res, err := c.client.Do(req)
	if ho.retryHook != nil {
		return ho.retryHook(req, res, c.client, err)
	}

	if ho.responseHook != nil {
		if err := ho.responseHook(req, res); err != nil {
			return nil, fmt.Errorf("failed to execute response hook: %w", err)
		}
	}

	return res, nil
}
