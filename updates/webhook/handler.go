package webhook

import (
	"io"
	"net/http"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/updates"
)

// Handler is a handler for telegram webhook requests.
type Handler struct {
	Handler updates.Handler
	Secret  string

	// OnError is called when telegram sends an invalid request.
	// The error has an Error type.
	OnError func(HTTPRequest, error)
}

func (wh Handler) onError(w HTTPResponse, r HTTPRequest, e Error) {
	if wh.OnError != nil {
		wh.OnError(r, e)
	}
	e.write(w)
}

// Handle handles the HTTP request.
//
//nolint:errcheck
func (wh *Handler) Handle(w HTTPResponse, r HTTPRequest) {
	if r.Method() != http.MethodPost {
		wh.onError(w, r, Error{
			Code: http.StatusMethodNotAllowed,
			Err:  "wrong http method (POST is required)",
		})
		return
	}

	if wh.Secret != "" &&
		wh.Secret != r.Header("X-Telegram-Bot-Api-Secret-Token") {
		wh.onError(w, r, Error{
			Code: http.StatusForbidden,
			Err:  "wrong secret token",
		})
		return
	}

	upd, _, err := api.DecodeJSON[tg.Update](r.Body())
	if err != nil {
		wh.onError(w, r, Error{
			Code: http.StatusBadRequest,
			Err:  err.Error(),
		})
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
func (wh *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wh.Handle(w, NewHTTPRequest(r))
}
