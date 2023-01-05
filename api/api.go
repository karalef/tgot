package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karalef/tgot/api/tg"
)

// New creates a new API instance.
// If apiURL or fileURL are empty, the Telegram defaults are used.
// If client is nil, the http.DefaultClient is used.
func New(token string, apiURL, fileURL string, client Client) (*API, error) {
	if token == "" {
		return nil, errors.New("no token provided")
	}
	if apiURL == "" {
		apiURL = tg.DefaultAPIURL
	}
	if fileURL == "" {
		fileURL = tg.DefaultFileURL
	}
	if client == nil {
		client = WrapStdHTTP(http.DefaultClient)
	}
	return &API{
		token:   token,
		apiURL:  apiURL,
		fileURL: fileURL,
		client:  client,
	}, nil
}

// NewDefault creates a new API instance with default values.
func NewDefault(token string) (*API, error) {
	return New(token, "", "", nil)
}

// API provides access to the Telegram Bot API.
type API struct {
	token   string
	apiURL  string
	fileURL string
	client  Client
}

// Request performs a request to the Bot API with background context,
// but doesn't parse the result.
func (a *API) Request(method string, d *Data) error {
	_, err := Request[Empty](a, method, d)
	return err
}

// Request performs a request to the Bot API with background context.
func Request[T any](a *API, method string, d *Data) (T, error) {
	return RequestContext[T](context.Background(), a, method, d)
}

// RequestContext performs a request to the Bot API.
func RequestContext[T any](ctx context.Context, a *API, method string, data *Data) (T, error) {
	ctype, reader := data.Data()

	var nilResult T
	u := a.apiURL + a.token + "/" + method
	body, err := a.client.Post(ctx, u, ctype, reader)
	switch {
	case err == nil:
	case errors.Is(err, context.Canceled):
		return nilResult, context.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return nilResult, context.DeadlineExceeded
	default:
		return nilResult, &HTTPError{makeError(method, data, err)}
	}
	defer body.Close()

	r, raw, err := DecodeJSON[tg.APIResponse[T]](body)
	if err != nil {
		return nilResult, &JSONError{
			baseError: makeError(method, data, err),
			Response:  raw,
		}
	}
	if r.APIError != nil {
		err = &Error{makeError(method, data, r.APIError)}
	}

	return r.Result, err
}

// DownloadFile downloads a file from the server.
func (a *API) DownloadFile(path string) (io.ReadCloser, error) {
	return a.client.Get(a.fileURL + a.token + "/" + path)
}

// DecodeJSON decodes reader into object or
// returns raw json data if error occured.
func DecodeJSON[T any](r io.Reader) (*T, []byte, error) {
	var v T
	dec := json.NewDecoder(r)
	err := dec.Decode(&v)
	if err == nil {
		return &v, nil, nil
	}
	b, _ := io.ReadAll(io.MultiReader(dec.Buffered(), r))
	return nil, b, err
}

// Empty type is used to avoid spending resources on unmarshaling.
type Empty struct{}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Empty) UnmarshalJSON([]byte) error { return nil }
