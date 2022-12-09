package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karalef/tgot/api/tg"
)

// New creates a new API instance and returns the getMe result if successful.
// If apiURL or fileURL are empty, the Telegram defaults are used.
// If client is nil, the http.DefaultClient is used.
func New(token string, apiURL, fileURL string, client *http.Client) (*API, error) {
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
		client = http.DefaultClient
	}
	return &API{
		token:   token,
		apiURL:  apiURL,
		fileURL: fileURL,
		client:  client,
	}, nil
}

// API provides access to the Telegram Bot API.
type API struct {
	token   string
	apiURL  string
	fileURL string
	client  *http.Client
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, reader)
	if err != nil {
		return nilResult, &HTTPError{makeError(method, data, err)}
	}
	if data != nil {
		req.Header.Set("Content-Type", ctype)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		switch e := errors.Unwrap(err); e {
		case context.Canceled, context.DeadlineExceeded:
			return nilResult, e
		default:
			return nilResult, &HTTPError{makeError(method, data, err)}
		}
	}
	defer resp.Body.Close()

	r, raw, err := DecodeJSON[tg.APIResponse[T]](resp.Body)
	if err != nil {
		return nilResult, &JSONError{
			baseError: makeError(method, data, err),
			Response:  raw,
		}
	}
	if r.APIError != nil {
		err = &Error{makeError(method, data, *r.APIError)}
	}

	return r.Result, err
}

// DownloadFile downloads a file from the server.
func (a *API) DownloadFile(path string) (io.ReadCloser, error) {
	resp, err := a.client.Get(a.fileURL + a.token + "/" + path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// GetMe returns basic information about the bot in form of a User object.
func (a *API) GetMe() (*tg.User, error) {
	return Request[*tg.User](a, "getMe", nil)
}

// LogOut method.
//
// Use this method to log out from the cloud Bot API server before launching the bot locally.
func (a *API) LogOut() error {
	return a.Request("logOut", nil)
}

// Close method.
//
// Use this method to close the bot instance before moving it from one local server to another.
func (a *API) Close() error {
	return a.Request("close", nil)
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
