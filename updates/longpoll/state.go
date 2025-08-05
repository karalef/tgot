package longpoll

import (
	"context"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Config contains the updates offset and the configuration for the long poller.
type Config struct {
	Offset  tg.ID
	Limit   uint
	Timeout uint
	Allowed []string
}

// NewState creates a new long poller state.
func NewState(a *api.API, cfg ...Config) *State {
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

// State is an optimized long polling interface. It does not allocate the
// api.Data object and marshals just the offset on each call, so it is not safe
// for concurrent use. State polls for updates using the getUpdates method and
// stores the offset.
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
// Calling the Poll method after Close causes panic.
func (s *State) Close() {
	s.d.Put()
	s.a = nil
}
