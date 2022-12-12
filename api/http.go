package api

import (
	"context"
	"io"
	"net/http"
)

// Client represents HTTP client.
type Client interface {
	Get(url string) (respBody io.ReadCloser, err error)

	// if the context is cancelled, it should return a context error or
	// an error that can be identified as a context error via errors.Is.
	Post(ctx context.Context, url string, contentType string, body io.Reader) (respBody io.ReadCloser, err error)
}

// WrapStdHTTP wraps std http.Client.
func WrapStdHTTP(c *http.Client) Client {
	return stdHTTPWrapper{c}
}

type stdHTTPWrapper struct {
	*http.Client
}

func (c stdHTTPWrapper) Get(url string) (io.ReadCloser, error) {
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

func (c stdHTTPWrapper) Post(ctx context.Context, url string, contentType string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
