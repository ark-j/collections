package httpx

import "net/http"

type (
	ResponseHook func(*http.Request, *http.Response) error
	RequestHook  func(*http.Request) error
	RetryHook    func(*http.Request, *http.Response, *http.Client, error) (*http.Response, error)
)

type HTTPOptions struct {
	headers      map[string]string
	queries      map[string]string
	responseHook ResponseHook
	requestHook  RequestHook
	retryHook    RetryHook
}

func NewHTTPOptions() *HTTPOptions {
	return &HTTPOptions{
		headers: make(map[string]string),
		queries: make(map[string]string),
	}
}

func (ho *HTTPOptions) Header(k, v string) *HTTPOptions {
	ho.headers[k] = v
	return ho
}

func (ho *HTTPOptions) Headers(hdrs map[string]string) *HTTPOptions {
	ho.headers = hdrs
	return ho
}

func (ho *HTTPOptions) Query(k, v string) *HTTPOptions {
	ho.queries[k] = v
	return ho
}

func (ho *HTTPOptions) Queries(queries map[string]string) *HTTPOptions {
	ho.queries = queries
	return ho
}

func (ho *HTTPOptions) RequestHook(hook RequestHook) *HTTPOptions {
	ho.requestHook = hook
	return ho
}

func (ho *HTTPOptions) ResponseHook(hook ResponseHook) *HTTPOptions {
	ho.responseHook = hook
	return ho
}

func (ho *HTTPOptions) RetryHook(hook RetryHook) *HTTPOptions {
	ho.retryHook = hook
	return ho
}
