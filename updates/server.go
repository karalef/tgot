package updates

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// NewWebhookServer creates new server for telegram webhooks.
// Ports currently supported for webhooks: 443, 80, 88, 8443.
func NewWebhookServer(addr string, cfg ServerConfig) *Server {
	if cfg.URL == "" || cfg.CertFile != "" && cfg.KeyFile == "" {
		return nil
	}
	mux := http.NewServeMux()
	return &Server{
		Serv: http.Server{
			Addr:    addr,
			Handler: mux,
		},
		Mux: mux,
		cfg: cfg,
	}
}

var _ Poller = &Server{}
var _ WebhookPoller = &Server{}

// Server is a server that handles telegram webhooks.
// It represents complete Poller.
// Mux can be used
type Server struct {
	// tls config will be automatically loaded if CertFile and KeyFile are specified
	Serv http.Server
	// can be used to set paths other than the one specified in config (it will be overwritten)
	Mux *http.ServeMux
	cfg ServerConfig
}

// ServerConfig contains webhook parameters.
type ServerConfig struct {
	Path     string // template for mux
	Filter   FilterFunc
	CertFile string // will be automatically opened and sent with setWebhook as "certificate"
	KeyFile  string

	// HTTPS URL to send updates to.
	URL string
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

// Close shuts down the server without interrupting any active connections.
func (s *Server) Close() {
	s.Serv.Shutdown(context.Background())
}

// Run starts webhook server.
func (s *Server) Run(a *api.API, h Handler, allowed []string) error {
	whhandler := func(upd *tg.Update) (string, *api.Data) {
		h(upd)
		return "", nil
	}
	return s.RunWH(a, whhandler, allowed)
}

// RunWH starts webhook server.
func (s *Server) RunWH(api *api.API, h WHHandler, allowed []string) error {
	if h == nil {
		panic("Webhooker: nil handler")
	}

	tls := s.cfg.CertFile != ""
	var cert *tg.InputFile
	if tls {
		certfile, err := os.ReadFile(s.cfg.CertFile)
		if err != nil {
			return err
		}
		cert = tg.FileBytes(filepath.Base(s.cfg.CertFile), certfile)
	}

	ok, err := SetWebhook(api, WebhookData{
		URL:            s.cfg.URL,
		Certificate:    cert,
		IPAddress:      s.cfg.IPAddress,
		MaxConnections: s.cfg.MaxConnections,
		AllowedUpdates: allowed,
		DropPending:    s.cfg.DropPending,
		SecretToken:    s.cfg.SecretToken,
	})
	if !ok { // because on successful result it returns an error with code 0
		return err
	}

	s.Mux.Handle(s.cfg.Path, &WebhookHandler{
		SecretToken: s.cfg.SecretToken,
		Handler:     h,
		Filter:      s.cfg.Filter,
	})

	if !tls {
		err = s.Serv.ListenAndServe()
	} else {
		err = s.Serv.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}