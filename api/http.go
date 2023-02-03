package api

import (
	"context"
	"io"
	"net/http"
)

// Client represents HTTP client.
type Client interface {
	Get(ctx context.Context, url string) (status int, respBody io.ReadCloser, err error)

	Post(ctx context.Context, url string, contentType string, body io.Reader) (status int, respBody io.ReadCloser, err error)
}

var defaultClient = WrapStdHTTP(http.DefaultClient)

// WrapStdHTTP wraps std http.Client.
func WrapStdHTTP(c *http.Client) Client {
	return stdHTTPWrapper{c}
}

type stdHTTPWrapper struct {
	*http.Client
}

func (c stdHTTPWrapper) Get(ctx context.Context, url string) (int, io.ReadCloser, error) {
	resp, err := c.Client.Get(url)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, resp.Body, err
}

func (c stdHTTPWrapper) Post(ctx context.Context, url string, contentType string, body io.Reader) (int, io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return 0, nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, resp.Body, nil
}
