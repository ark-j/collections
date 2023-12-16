package reqwest

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
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
// TODO: enable http trace as optional
type Reqwest struct {
	client          *http.Client
	enableHTTPTrace bool
}

func New() *Reqwest {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		panic(err)
	}
	return (&Reqwest{}).SetTransport(nil).SetCookieJar(jar)
}

func (r *Reqwest) SetTransport(t *http.Transport) *Reqwest {
	if t != nil {
		r.client.Transport = t
	} else {
		r.client.Transport = defaultTransport
	}
	return r
}

func (r *Reqwest) SetCookieJar(jar http.CookieJar) *Reqwest {
	r.client.Jar = jar
	return r
}

func (r *Reqwest) GetClient() *http.Client {
	return r.client
}

func (r *Reqwest) Get(ctx context.Context, uri string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	return r.client.Do(req)
}

func (r *Reqwest) Head(ctx context.Context, uri string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	return r.client.Do(req)
}

func (r *Reqwest) Post(ctx context.Context, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	return r.client.Do(req)
}

func (r *Reqwest) Put(ctx context.Context, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, body)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	return r.client.Do(req)
}
