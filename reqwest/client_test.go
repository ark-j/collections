package reqwest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type Payload struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func TestClient(t *testing.T) {
	cases := []struct {
		name string
		want any
		exec func(t *testing.T, want any, uri string)
	}{
		{
			name: "test-http-get",
			want: Payload{Name: "GET", ID: 1},
			exec: func(t *testing.T, want any, uri string) {
				res, err := New(false).Get(context.Background(), uri)
				if noerr(t, err) {
					return
				}
				equals(t, http.StatusOK, res.StatusCode)
				defer res.Body.Close()

				got := Payload{}
				err = json.NewDecoder(res.Body).Decode(&got)
				noerr(t, err)
				equals(t, want, got)
			},
		},
		{
			name: "test-http-post",
			want: Payload{Name: "POST", ID: 2},
			exec: func(t *testing.T, want any, uri string) {
				var buf bytes.Buffer
				// nolint
				json.NewEncoder(&buf).Encode(want)
				res, err := New(false).Post(context.Background(), uri, &buf)
				if noerr(t, err) {
					return
				}
				equals(t, http.StatusCreated, res.StatusCode)
				defer res.Body.Close()

				var got Payload
				// nolint
				json.NewDecoder(res.Body).Decode(&got)
				equals(t, want, got)
			},
		},
		{
			name: "test-http-put",
			want: http.StatusAccepted,
			exec: func(t *testing.T, want any, uri string) {
				var buf bytes.Buffer
				body := Payload{Name: "PUT", ID: 3}
				// nolint
				json.NewEncoder(&buf).Encode(body)
				res, err := New(false).Put(context.Background(), uri, &buf)
				if noerr(t, err) {
					return
				}
				equals(t, want, res.StatusCode)
				defer res.Body.Close()
			},
		},
		{
			name: "test-http-delete",
			want: http.StatusNoContent,
			exec: func(t *testing.T, want any, uri string) {
				res, err := New(false).Delete(context.Background(), uri)
				if noerr(t, err) {
					return
				}
				equals(t, want, res.StatusCode)
				res.Body.Close()
			},
		},
		{
			name: "test-http-head",
			want: http.Header{"User_id": []string{"1111"}},
			exec: func(t *testing.T, want any, uri string) {
				res, err := New(false).Head(context.Background(), uri)
				if noerr(t, err) {
					return
				}
				defer res.Body.Close()
				equals(t, want, res.Header)
			},
		},
	}

	ts := mockHTTPServer()
	t.Cleanup(ts.Close)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.exec(t, tt.want, ts.URL)
		})
	}
}

// mockHTTPServer for testing http client
func mockHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			body := &Payload{Name: "GET", ID: 1}
			// nolint
			json.NewEncoder(w).Encode(body)
		case http.MethodPut:
			var body Payload
			// nolint
			json.NewDecoder(r.Body).Decode(&body)
			if body.ID != 3 {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			// nolint
			json.NewEncoder(w).Encode(&body)
		case http.MethodPost:
			var body Payload
			// nolint
			json.NewDecoder(r.Body).Decode(&body)
			w.WriteHeader(http.StatusCreated)
			// nolint
			json.NewEncoder(w).Encode(&body)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodOptions:
		case http.MethodHead:
			w.Header().Set("user_id", "1111")
		}
	}))
}

// Custom Transport for roundTripper
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// new test client
func NewTestClient(fn RoundTripFunc) *Reqwest {
	return New(false).SetTransport(RoundTripFunc(fn))
}

// simple client test
func TestClientMisc(t *testing.T) {
	cases := []struct {
		name         string
		want         any
		options      []Options
		roundTripper func(want any) RoundTripFunc
	}{
		{
			name:    "test-http-query",
			want:    url.Values{"user": []string{"1111"}},
			options: []Options{WithQuries(map[string]string{"user": "1111"})},
			roundTripper: func(want any) RoundTripFunc {
				return func(req *http.Request) *http.Response {
					equals(t, req.URL.Query(), want)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString("OK")),
						Header:     make(http.Header),
					}
				}
			},
		},
		{
			name: "test-http-headers",
			want: http.Header{
				"User":       []string{"1111"},
				"User-Agent": []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
			},
			options: []Options{WithHeaders(map[string]string{"User": "1111"})},
			roundTripper: func(want any) RoundTripFunc {
				return func(req *http.Request) *http.Response {
					equals(t, req.Header, want)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString("OK")),
						Header:     make(http.Header),
					}
				}
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := NewTestClient(tt.roundTripper(tt.want))

			res, err := c.Get(context.Background(), "https://example.com", tt.options...)
			if noerr(t, err) {
				return
			}
			defer res.Body.Close()
		})
	}
}

// helper for equality
func equals(t testing.TB, got, want any) bool {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("wanted %v got %v", want, got)
		return false
	}
	return true
}

// helper utility for noerr
func noerr(t testing.TB, err error) bool {
	t.Helper()
	if err != nil {
		t.Errorf("required no err but got err:%v", err)
		return false
	}
	return true
}
