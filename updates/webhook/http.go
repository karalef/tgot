package webhook

import (
	"io"
	"net/http"
)

// HTTPRequest represents a server-side HTTP request from client.
type HTTPRequest interface {
	// Method returns the HTTP method of the request.
	Method() string

	// Header returns the HTTP header by key.
	Header(key string) string

	// Body returns the request body.
	Body() io.Reader
}

// HTTPResponse represents a server-side HTTP response to client.
type HTTPResponse http.ResponseWriter

// Error represents a telegram-side error.
type Error struct {
	Code int
	Err  string
}

// Error implements error interface.
func (e Error) Error() string { return e.Err }

//nolint:errcheck
func (e Error) write(w HTTPResponse) {
	w.WriteHeader(e.Code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"error\":\"" + e.Err + "\"}"))
}

// NewHTTPRequest creates a new request from the given std HTTP request.
func NewHTTPRequest(r *http.Request) HTTPRequest { return stdRequest{r} }

type stdRequest struct{ *http.Request }

func (r stdRequest) Method() string           { return r.Request.Method }
func (r stdRequest) Header(key string) string { return r.Request.Header.Get(key) }
func (r stdRequest) Body() io.Reader          { return r.Request.Body }
