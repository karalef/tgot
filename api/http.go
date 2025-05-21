package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// HTTP represents HTTP client that can do GET and POST requests.
type HTTP interface {
	// Get performs a GET request with context to the specified URL.
	// It returns status code and response body.
	Get(ctx context.Context, url string) (int, io.ReadCloser, error)

	// Post performs a POST request with context to the specified URL using the body.
	// It returns status code and response body.
	Post(ctx context.Context, url, ct string, body io.Reader) (int, io.ReadCloser, error)
}

// HTTPError represents HTTP error.
type HTTPError struct {
	Status int
	Err    error
	URL    string
}

func (e *HTTPError) Unwrap() error { return e.Err }

func (e *HTTPError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("http %s (%s)", e.URL, httpStatus(e.Status))
	}
	return fmt.Sprintf("http %s (%s): %s", e.URL, httpStatus(e.Status), e.Err.Error())
}

func httpStatus(code int) string {
	return strconv.Itoa(code) + " " + http.StatusText(code)
}

// NewHTTP returns HTTP implementation that uses std http.Client.
func NewHTTP(c *http.Client) HTTP { return (*stdHTTP)(c) }

type stdHTTP http.Client

func (h *stdHTTP) Get(ctx context.Context, url string) (int, io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}
	resp, err := (*http.Client)(h).Do(req)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, resp.Body, nil
}

func (h *stdHTTP) Post(ctx context.Context, url, ct string, body io.Reader) (int, io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", ct)
	resp, err := (*http.Client)(h).Do(req)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, resp.Body, nil
}
