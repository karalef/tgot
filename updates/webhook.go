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

// WrapHandler wraps Handler with WHHandler.
func WrapHandler(h Handler) WHHandler {
	return func(upd *tg.Update) (string, *api.Data, error) {
		return "", nil, h(upd)
	}
}

// WebhookHandler is a handler for telegram webhook requests.
// It implements std http.Handler.
type WebhookHandler struct {
	SecretToken string
	Handler     WHHandler
}

func writeErr(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"error\":\"" + err + "\"}"))
}

func (wh *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "wrong http method (POST is required)")
		return
	}

	if wh.SecretToken != "" &&
		wh.SecretToken != r.Header.Get("X-Telegram-Bot-Api-Secret-Token") {
		writeErr(w, http.StatusUnauthorized, "wrong secret token")
		return
	}

	upd, _, err := api.DecodeJSON[tg.Update](r.Body)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	meth, data, err := wh.Handler(upd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	if meth == "" {
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
