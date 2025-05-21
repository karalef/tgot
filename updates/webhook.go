package updates

import (
	"io"
	"net/http"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// WHHandler represents webhook handler function type.
// If the method is not empty, the request to the api will be written in response to the webhook.
type WHHandler func(*tg.Update) (method string, data *api.Data, err error)

// WrapWebhook wraps Handler with WHHandler.
func WrapWebhook(h Handler) WHHandler {
	return func(upd *tg.Update) (string, *api.Data, error) {
		return "", nil, h(upd)
	}
}

// WebhookHandler is a handler for telegram webhook requests.
// It implements std http.Handler.
type WebhookHandler struct {
	Handler     WHHandler
	SecretToken string
}

func writeErr(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"error\":\"" + err + "\"}"))
}

// Handle handles the HTTP request.
func (wh *WebhookHandler) Handle(w Response, r Request) {
	if r.Method() != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "wrong http method (POST is required)")
		return
	}

	if wh.SecretToken != "" &&
		wh.SecretToken != r.Header("X-Telegram-Bot-Api-Secret-Token") {
		writeErr(w, http.StatusUnauthorized, "wrong secret token")
		return
	}

	upd, _, err := api.DecodeJSON[tg.Update](r.Body())
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	meth, data, err := wh.Handler(upd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if meth == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if data == nil {
		data = api.NewData()
	}
	data.Set("method", meth)
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

// Response represents a server-side HTTP response to client.
type Response http.ResponseWriter

// NewRequest creates a new request from the given std HTTP request.
func NewRequest(r *http.Request) Request { return stdRequest{r} }

type stdRequest struct{ *http.Request }

func (r stdRequest) Method() string           { return r.Request.Method }
func (r stdRequest) Header(key string) string { return r.Request.Header.Get(key) }
func (r stdRequest) Body() io.Reader          { return r.Request.Body }
