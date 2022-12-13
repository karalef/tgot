package updates

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// WHHandler represents webhook handler function type.
// If the method is not empty, the request to the api will be written in response to the webhook.
type WHHandler func(*tg.Update) (method string, data *api.Data)

// WebhookPoller represents a webhook server that can send api requests in response to webhook requests.
type WebhookPoller interface {
	Poller
	RunWH(api *api.API, handler WHHandler, allowed []string) error
}

// WrapHandler wraps Handler with WHHandler.
func WrapHandler(h Handler) WHHandler {
	return func(upd *tg.Update) (string, *api.Data) {
		h(upd)
		return "", nil
	}
}

// WebhookHandler is a handler for telegram webhook requests.
type WebhookHandler struct {
	SecretToken string
	Handler     WHHandler
	Filter      FilterFunc
}

func writeErr(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"error\":\"" + err + "\"}"))
}

func (wh *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusBadRequest, "wrong http method (POST is required)")
		return
	}

	if wh.SecretToken != "" &&
		wh.SecretToken != r.Header.Get("X-Telegram-Bot-Api-Secret-Token") {
		writeErr(w, http.StatusUnauthorized, "wrong secret token")
		return
	}

	var upd tg.Update
	err := json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

	if wh.Filter != nil && !wh.Filter(&upd) {
		return
	}

	meth, data := wh.Handler(&upd)
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
