package longpoll

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/updates"
)

// NewLongPoller creates a new long poller.
func NewLongPoller(api *api.API) *LongPoller { return &LongPoller{api: api} }

// LongPoller uses long polling to listen to the server for updates via the getUpdates method.
type LongPoller struct {
	api *api.API
	run atomic.Bool
	wg  sync.WaitGroup
}

// Run starts long polling. The passed context controlls only the long poller.
// It can be called multiple times but only after stopping via context
// cancellation.
func (lp *LongPoller) Run(ctx context.Context, h updates.Handler, cfg ...Config) error {
	if ctx == nil {
		panic("LongPoller: nil context")
	}
	if h == nil {
		panic("LongPoller: nil handler")
	}
	if lp.run.CompareAndSwap(false, true) {
		panic("LongPoller: already running")
	}

	s := NewState(lp.api, cfg...)
	defer lp.run.Store(false)
	defer lp.wg.Wait()
	defer s.Close()

	for {
		upds, err := s.Poll(ctx)
		if err != nil {
			return err
		}
		for i := range upds {
			go lp.handle(ctx, h, &upds[i])
		}
	}
}

func (lp *LongPoller) handle(ctx context.Context, h updates.Handler, upd *tg.Update) {
	lp.wg.Add(1)
	defer lp.wg.Done()

	if resp := h.Handle(upd); resp != nil {
		_ = lp.api.Request(ctx, resp.Method(), resp.Data())
	}
}
