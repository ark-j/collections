package reqwest

type HTTPOptions struct {
	headers map[string]string
	queries map[string]string
}

type Options func(o *HTTPOptions)

// WithQueries is Options to provide queries
func WithQuries(q map[string]string) Options {
	return func(o *HTTPOptions) {
		o.queries = q
	}
}

// WithHeaders is option to provide headers
func WithHeaders(h map[string]string) Options {
	return func(o *HTTPOptions) {
		o.headers = h
	}
}
