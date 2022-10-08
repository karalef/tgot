package updates

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

// NewWebhook creates new server for telegram webhooks.
// Ports currently supported for webhooks: 443, 80, 88, 8443.
func NewWebhook(port int, filter Filter, cfg WebhookConfig) *Webhooker {
	if cfg.CertFile != "" && cfg.KeyFile == "" || cfg.URL == "" {
		return nil
	}
	wh := &Webhooker{cfg: cfg, Filter: filter}
	wh.serv = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: http.HandlerFunc(wh.handle),
	}
	return wh
}

// Webhooker receives incoming updates via an outgoing webhook.
type Webhooker struct {
	serv    *http.Server
	cfg     WebhookConfig
	handler WHHandler

	Filter Filter
}

// WebhookConfig contains webhook parameters.
type WebhookConfig struct {
	// HTTPS URL to send updates to.
	URL string

	CertFile string
	KeyFile  string

	// The fixed IP address which will be used to send webhook requests instead of the IP address resolved through DNS.
	IPAddress string
	// The maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery, 1-100.
	// Defaults to 40.
	MaxConnections int
	// Pass True to drop all pending updates.
	DropPending bool
	// A secret token to be sent in a header “X-Telegram-Bot-Api-Secret-Token” in every webhook request, 1-256 characters.
	// Only characters A-Z, a-z, 0-9, _ and - are allowed.
	// The header is useful to ensure that the request comes from a webhook set by you.
	SecretToken string
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (wh *Webhooker) Shutdown() {
	wh.serv.Shutdown(context.Background())
}

// Close immediately stops the server.
func (wh *Webhooker) Close() {
	wh.serv.Close()
}

// Run starts webhook server.
func (wh *Webhooker) Run(a *api.API, h Handler, allowed []string) error {
	whhandler := func(upd *tg.Update) (string, api.Data) {
		h(upd)
		return "", api.Data{}
	}
	return wh.RunWH(a, whhandler, allowed)
}

// RunWH starts webhook server.
func (wh *Webhooker) RunWH(api *api.API, h WHHandler, allowed []string) error {
	if h == nil {
		panic("Webhooker: nil handler")
	}

	tls := wh.cfg.CertFile != ""
	var cert *tg.InputFile
	if tls {
		certfile, err := os.ReadFile(wh.cfg.CertFile)
		if err != nil {
			return err
		}
		cert = tg.FileBytes(filepath.Base(wh.cfg.CertFile), certfile)
	}

	err := SetWebhook(api, WebhookData{
		URL:            wh.cfg.URL,
		Certificate:    cert,
		IPAddress:      wh.cfg.IPAddress,
		MaxConnections: wh.cfg.MaxConnections,
		AllowedUpdates: allowed,
		DropPending:    wh.cfg.DropPending,
		SecretToken:    wh.cfg.SecretToken,
	})
	if err != nil {
		return err
	}

	wh.handler = h
	if !tls {
		err = wh.serv.ListenAndServe()
	} else {
		err = wh.serv.ListenAndServeTLS(wh.cfg.CertFile, wh.cfg.KeyFile)
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	wh.Shutdown()
	return err
}

func writeErr(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"error\":\"" + err + "\"}"))
}

func (wh *Webhooker) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, "wrong http method (POST is required)")
		return
	}

	if wh.cfg.SecretToken != "" &&
		wh.cfg.SecretToken != r.Header.Get("X-Telegram-Bot-Api-Secret-Token") {
		writeErr(w, "wrong secret token")
		return
	}

	var upd tg.Update
	err := json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		writeErr(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

	if wh.Filter != nil && !wh.Filter(&upd) {
		return
	}

	meth, data := wh.handler(&upd)
	if meth == "" {
		return
	}
	data.Set("method", meth)
	ctype, reader := data.Data()
	w.Header().Set("Content-Type", ctype)
	io.Copy(w, reader)
}

// WebhookData contains parameters for setWebhook method.
type WebhookData struct {
	URL            string
	Certificate    *tg.InputFile
	IPAddress      string
	MaxConnections int
	AllowedUpdates []string
	DropPending    bool
	SecretToken    string
}

// SetWebhook specifies a webhook URL.
// Use this method to specify a URL and receive incoming updates via an outgoing webhook.
func SetWebhook(a *api.API, s WebhookData) error {
	d := api.NewData().Set("url", s.URL)
	if s.Certificate != nil {
		d.Files = []api.File{{
			Field:      "certificate",
			Inputtable: s.Certificate,
		}}
	}
	d.Set("ip_address", s.IPAddress)
	d.SetInt("max_connections", s.MaxConnections)
	d.SetJSON("allowed_updates", s.AllowedUpdates)
	d.SetBool("drop_pending_updates", s.DropPending)
	d.Set("secret_token", s.SecretToken)
	return a.Request("setWebhook", d)
}

// DeleteWebhook removes webhook integration if you decide to switch back to getUpdates.
func DeleteWebhook(a *api.API, dropPending bool) error {
	d := api.NewData().SetBool("drop_pending_updates", dropPending)
	return a.Request("deleteWebhook", d)
}

// GetWebhookInfo returns current webhook status.
func GetWebhookInfo(a *api.API) (*tg.WebhookInfo, error) {
	return api.Request[*tg.WebhookInfo](a, "getWebhookInfo")
}
