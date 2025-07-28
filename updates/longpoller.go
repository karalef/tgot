package updates

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
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
func (lp *LongPoller) Run(ctx context.Context, h Handler, cfg ...LongPollConfig) error {
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

func (lp *LongPoller) handle(ctx context.Context, h Handler, upd *tg.Update) {
	lp.wg.Add(1)
	defer lp.wg.Done()

	if resp := h.Handle(upd); resp != nil {
		_ = lp.api.Request(ctx, resp.Method(), resp.Data())
	}
}

// LongPollConfig is the configuration for the long poller.
type LongPollConfig struct {
	Offset  tg.ID
	Limit   uint
	Timeout uint
	Allowed []string
}

// NewState creates a new long poller state.
func NewState(a *api.API, cfg ...LongPollConfig) *State {
	s := &State{
		a: a,
		d: api.NewData(),
	}
	if len(cfg) > 0 {
		s.o = cfg[0].Offset
		s.d.
			SetUint("limit", cfg[0].Limit).
			SetUint("timeout", cfg[0].Timeout).
			SetJSON("allowed", cfg[0].Allowed)
	}
	return s
}

// State polls for updates using the getUpdates method and stores the offset.
// It is not safe for concurrent use.
type State struct {
	a *api.API
	d *api.Data
	o tg.ID
}

// Poll polls for updates from the server and updates the offset.
func (s *State) Poll(ctx context.Context) ([]tg.Update, error) {
	s.d.SetID("offset", s.o)
	upds, err := api.Request[[]tg.Update](ctx, s.a, "getUpdates", s.d)
	if err != nil {
		return nil, err
	}
	if len(upds) > 0 {
		s.o = upds[len(upds)-1].ID + 1
	}
	return upds, nil
}

// Close releases the resources used by the state.
// Calling the Receive method after Close causes panic.
func (s *State) Close() {
	s.d.Put()
	s.a = nil
}
