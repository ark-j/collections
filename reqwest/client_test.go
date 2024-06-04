package reqwest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"
)

// Custom Transport for roundTripper
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestClient(t *testing.T) {
	t.Run("client-test", simpleClientTest)
}

func testMethods(t *testing.T) {
	type Payload struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
	cases := []struct {
		name   string
		method string
		want   Payload
		got    Payload
	}{
		{
			name:   "http-get",
			method: http.MethodGet,
			want:   Payload{},
			got:    Payload{},
		},
		{
			name:   "http-post",
			method: http.MethodPost,
			want:   Payload{},
			got:    Payload{},
		},
		{
			name:   "http-put",
			method: http.MethodPut,
			want:   Payload{},
			got:    Payload{},
		},
		{
			name:   "http-delete",
			method: http.MethodDelete,
			want:   Payload{},
			got:    Payload{},
		},
		{
			name:   "http-options",
			method: http.MethodOptions,
			want:   Payload{},
			got:    Payload{},
		},
		{
			name:   "http-head",
			method: http.MethodHead,
			want:   Payload{},
			got:    Payload{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

// simple client test
func simpleClientTest(t *testing.T) {
	c := NewTestClient(func(req *http.Request) *http.Response {
		equals(t, req.URL.String(), "https://example.com")
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("OK")),
			Header:     make(http.Header),
		}
	})

	res, err := c.Get(context.Background(), "https://example.com", nil, nil)
	if err != nil {
		t.Errorf("err while getting response: %v", err)
		return
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	equals(t, b, []byte("OK"))
}

// new test client
func NewTestClient(fn RoundTripFunc) *Reqwest {
	return New(false).SetTransport(RoundTripFunc(fn))
}

// helper for equality
func equals(t testing.TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wanted %v got %v", want, got)
	}
}
