package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karalef/tgot/api/tg"
)

// DefaultAPIURL is a default url for telegram api.
const DefaultAPIURL = "https://api.telegram.org/bot"

// DefaultFileURL is a default url for downloading files.
const DefaultFileURL = "https://api.telegram.org/file/bot"

// New creates a new API instance.
// If apiURL or fileURL are empty, the Telegram defaults are used.
// If client is nil, the http.DefaultClient is used.
func New(token string, apiURL, fileURL string, client *http.Client) (*API, error) {
	if token == "" {
		return nil, errors.New("no token provided")
	}
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}
	if fileURL == "" {
		fileURL = DefaultFileURL
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &API{
		token:      token,
		apiURL:     apiURL,
		fileURL:    fileURL,
		botAPIURL:  apiURL + token + "/",
		botFileURL: fileURL + token + "/",
		client:     client,
	}, nil
}

// NewDefault creates a new API instance with default values.
func NewDefault(token string) (*API, error) {
	return New(token, "", "", nil)
}

// API provides access to the Telegram Bot API.
type API struct {
	token      string
	apiURL     string
	fileURL    string
	botAPIURL  string
	botFileURL string
	client     *http.Client
}

func (a API) methodURL(method string) string { return a.botAPIURL + method }
func (a API) pathURL(filepath string) string { return a.botFileURL + filepath }

// RequestContext performs a request to the Bot API but doesn't parse the result.
func (a *API) Request(ctx context.Context, method string, d *Data) error {
	_, err := Request[Empty](ctx, a, method, d)
	return err
}

// Request performs a request to the Bot API.
func Request[T any](ctx context.Context, a *API, method string, data *Data) (result T, err error) {
	ctype, reader := data.Data()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.methodURL(method), reader)
	if err != nil {
		return result, err
	}
	req.Header.Set("Content-Type", ctype)

	resp, err := a.client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	r, raw, err := DecodeJSON[tg.Response[T]](resp.Body)
	if err != nil {
		return result, &JSONError{
			baseError: makeError(method, data, err),
			Status:    resp.StatusCode,
			Response:  raw,
		}
	}
	if r.Error != nil {
		err = &Error{makeError(method, data, r.Error)}
	}

	return r.Result, err
}

// DownloadFile downloads a file from the server.
func (a *API) DownloadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.pathURL(path), nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, &DownloadError{
			Status: resp.StatusCode,
			Path:   path,
			Err:    err,
		}
	}
	return resp.Body, nil
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
func (Empty) UnmarshalJSON([]byte) error { return nil }
