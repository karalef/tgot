package webhook

import (
	"context"
	"errors"
	"net/http"

	"github.com/karalef/tgot/updates"
)

// NewServer creates new server for telegram webhooks.
// Ports currently supported for webhooks: 443, 80, 88, 8443.
func NewServer(addr string, cfg Config) (*Server, error) {
	if cfg.CertFile != "" && cfg.KeyFile == "" {
		return nil, errors.New("certificate file without key file")
	}
	mux := http.NewServeMux()
	return &Server{
		Serv: http.Server{
			Addr:    addr,
			Handler: mux,
		},
		Mux: mux,
		cfg: cfg,
	}, nil
}

// Server is a server that handles telegram webhooks.
// It uses std HTTP server implementation.
//
// Mux can be used for templates other than specified in path.
type Server struct {
	// tls config will be automatically loaded if CertFile and KeyFile are specified
	Serv http.Server
	// can be used to set paths other than the one specified in config (it will be overwritten)
	Mux *http.ServeMux
	cfg Config
}

// Config contains webhook parameters.
type Config struct {
	Path     string // template for mux
	CertFile string
	KeyFile  string
	Secret   string
}

// ListenAndServe starts webhook server.
func (s *Server) ListenAndServe(ctx context.Context, h updates.Handler) (err error) {
	if h == nil {
		panic("WebhookServer: nil handler")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.Mux.Handle(s.cfg.Path, &Handler{
		Handler: h,
		Secret:  s.cfg.Secret,
	})

	closed := make(chan error)
	go func() {
		<-ctx.Done()
		closed <- s.Serv.Shutdown(context.Background())
		close(closed)
	}()

	if s.cfg.CertFile == "" {
		err = s.Serv.ListenAndServe()
	} else {
		err = s.Serv.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	if errClosed := <-closed; errClosed != nil && err == nil {
		err = errClosed
	}
	return err
}
