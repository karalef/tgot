// package api provides the Telegram Bot API client.
// It allow to implement any API method.
// All the object types are defined in the underlying packages.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/karalef/tgot/api/tg"
)

// DefaultAPIURL is a default url for telegram api.
const DefaultAPIURL = "https://api.telegram.org/bot"

// DefaultFileURL is a default url for downloading files.
const DefaultFileURL = "https://api.telegram.org/file/bot"

func tokenURL(base, token string) string { return base + token + "/" }

// Config contains API parameters.
type Config struct {
	APIURL  string // default: DefaultAPIURL
	FileURL string // default: DefaultFileURL
	Client  HTTP   // default: http.DefaultClient
}

// New creates a new API instance.
func New(token string, cfg Config) (*API, error) {
	if token == "" {
		return nil, errors.New("no token provided")
	}
	if cfg.APIURL == "" {
		cfg.APIURL = DefaultAPIURL
	}
	if cfg.FileURL == "" {
		cfg.FileURL = DefaultFileURL
	}
	if cfg.Client == nil {
		cfg.Client = NewHTTP(http.DefaultClient)
	}
	return &API{
		token:   token,
		apiURL:  tokenURL(cfg.APIURL, token),
		fileURL: tokenURL(cfg.FileURL, token),
		client:  cfg.Client,
	}, nil
}

// API provides access to the Telegram Bot API.
type API struct {
	token   string
	apiURL  string
	fileURL string
	client  HTTP
}

func (a API) methodURL(method string) string { return a.apiURL + method }
func (a API) pathURL(filepath string) string { return a.fileURL + filepath }

func (a API) get(ctx context.Context, url string) (io.ReadCloser, *HTTPError) {
	code, body, err := a.client.Get(ctx, url)
	if err != nil || code != http.StatusOK {
		if body != nil {
			body.Close()
		}
		return nil, &HTTPError{
			Status: code,
			Err:    err,
			URL:    url,
		}
	}
	return body, nil
}

// RequestContext performs a request to the Bot API but doesn't parse the result.
func (a *API) Request(ctx context.Context, method string, d *Data) error {
	_, err := Request[Empty](ctx, a, method, d)
	return err
}

// Request performs a request to the Bot API.
func Request[T any](ctx context.Context, a *API, method string, data *Data) (result T, err error) {
	ctype, reader := data.Data()
	url := a.methodURL(method)
	code, body, err := a.client.Post(ctx, url, ctype, reader)
	if err != nil {
		return result, &HTTPError{
			Status: code,
			Err:    err,
			URL:    strings.Replace(url, a.token, "...", 1),
		}
	}
	defer body.Close()

	r, raw, err := DecodeJSON[tg.Response[T]](body)
	if err != nil {
		return result, &JSONError{
			baseError: makeError(method, data, err),
			Status:    code,
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
	body, err := a.get(ctx, a.pathURL(path))
	if err != nil {
		err.URL = strings.Replace(err.URL, a.token, "...", 1)
	}
	return body, err
}

func makeError[T error](method string, d *Data, err T) (e baseError[T]) {
	e.Method, e.Err = method, err
	if d != nil {
		e.Params = make(map[string]string, len(d.Params))
		for k, v := range d.Params {
			e.Params[k] = v
		}
		e.Files = make(map[string]string, len(d.Upload))
		for field, f := range d.Upload {
			s, _ := f.FileData()
			e.Files[field] = s
		}
	}
	return e
}

type baseError[T error] struct {
	Method string
	Params map[string]string
	Files  map[string]string
	Err    T
}

func (e baseError[T]) Unwrap() error { return e.Err }

func (e baseError[T]) Error() string {
	return fmt.Sprintf("%s\n%s %s", e.Err.Error(), e.Method, e.formatData())
}

func (e baseError[T]) formatData() string {
	var sb strings.Builder
	for k, v := range e.Params {
		sb.WriteString(k + "=" + v + " ")
	}
	for field, f := range e.Files {
		sb.WriteByte('\n')
		sb.WriteString("[file] " + field + ": " + f)
	}
	return sb.String()
}

// Error represents a telegram api error and also contains method and data.
type Error struct{ baseError[*tg.Error] }

// Is implements errors.Is interface.
func (e *Error) Is(err error) bool {
	if tge, ok := err.(*tg.Error); ok {
		return e.Err.Code == tge.Code
	}
	return false
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

// JSONError represents JSON error.
type JSONError struct {
	baseError[error]
	Status   int
	Response []byte
}

func (e *JSONError) Error() string {
	return fmt.Sprintf("%s\n%s %s %s", e.Err.Error(), httpStatus(e.Status), e.Method, e.formatData())
}

// Empty type is used to avoid spending resources on unmarshaling.
type Empty struct{}

// UnmarshalJSON implements json.Unmarshaler.
func (Empty) UnmarshalJSON([]byte) error { return nil }
