package updates

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Encoding is an interface for encoding.
type Encoding interface {
	EncodeToString([]byte) string
	DecodedLen(int) int
}

// GenerateSecret generates random secret token of specified length.
// It uses the provided encoding (or base64 without padding by default) to get
// the string representation of the secret.
// If length is 0 or bigger than 128, it will be set to 64.
func GenerateSecret(enc Encoding, length uint8) string {
	if enc == nil {
		enc = base64.RawStdEncoding
	}
	if length == 0 || length > 128 {
		length = 64
	}
	buf := make([]byte, enc.DecodedLen(int(length)))
	_, err := rand.Read(buf)
	if err != nil {
		panic("error while generating secret: " + err.Error())
	}
	return enc.EncodeToString(buf)
}

// WebhookHandler is a handler for telegram webhook requests.
type WebhookHandler struct {
	Handler     Handler
	SecretToken string
}

//nolint:errcheck
func writeErr(w HTTPResponse, code int, err string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"error\":\"" + err + "\"}"))
}

// Handle handles the HTTP request.
//
//nolint:errcheck
func (wh *WebhookHandler) Handle(w HTTPResponse, r Request) {
	if r.Method() != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "wrong http method (POST is required)")
		return
	}

	if wh.SecretToken != "" &&
		wh.SecretToken != r.Header("X-Telegram-Bot-Api-Secret-Token") {
		writeErr(w, http.StatusForbidden, "wrong secret token")
		return
	}

	upd, _, err := api.DecodeJSON[tg.Update](r.Body())
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := wh.Handler.Handle(upd)
	if resp == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	data := resp.Data()
	if data == nil {
		data = api.NewData()
	}
	data.Set("method", resp.Method())
	ctype, reader := data.Data()
	w.Header().Set("Content-Type", ctype)
	io.Copy(w, reader)
}

// ServeHTTP implements std http.Handler.
func (wh *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.Handle(w, NewRequest(r))
}

// Request represents a server-side HTTP request from client.
type Request interface {
	// Method returns the HTTP method of the request.
	Method() string

	// Header returns the HTTP header by key.
	Header(key string) string

	// Body returns the request body.
	Body() io.Reader
}

// HTTPResponse represents a server-side HTTP response to client.
type HTTPResponse http.ResponseWriter

// NewRequest creates a new request from the given std HTTP request.
func NewRequest(r *http.Request) Request { return stdRequest{r} }

type stdRequest struct{ *http.Request }

func (r stdRequest) Method() string           { return r.Request.Method }
func (r stdRequest) Header(key string) string { return r.Request.Header.Get(key) }
func (r stdRequest) Body() io.Reader          { return r.Request.Body }
