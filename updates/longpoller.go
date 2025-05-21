package updates

import (
	"context"
	"errors"
	"sync"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// NewLongPoller creates a LongPoller instance.
//
// limit must be in the range of 1-100 (else it will be set to server's default).
func NewLongPoller(timeout, limit uint, offset int) *LongPoller {
	if limit > 100 {
		limit = 0
	}
	lp := LongPoller{
		timeout: timeout,
		limit:   limit,
		offset:  offset,
	}
	return &lp
}

// LongPoller represents complete Poller that polls the server for updates via the getUpdates method.
type LongPoller struct {
	timeout, limit uint
	offset         int

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// Close stops the long poller and waits for all active handlers to complete.
// It panics if the poller is not running.
func (lp *LongPoller) Close() {
	lp.cancel()
	lp.wg.Wait()
}

// Run starts long polling. The passed context controlls only the long poller
// but not the handlers.
func (lp *LongPoller) Run(ctx context.Context, b *tgot.Bot) error {
	if b == nil {
		panic("LongPoller: nil bot")
	}

	ctx, lp.cancel = context.WithCancel(ctx)
	defer lp.wg.Wait()

	d := api.NewData()
	d.SetUint("limit", lp.limit)
	d.SetUint("timeout", lp.timeout)
	d.SetJSON("allowed", b.Allowed())
	defer d.Put()

	for a := b.API(); ; {
		d.SetInt("offset", lp.offset)
		upds, err := api.Request[[]tg.Update](ctx, a, "getUpdates", d)
		switch {
		case err == nil:
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return b.Err()
		default:
			return err
		}
		if len(upds) > 0 {
			lp.offset = upds[len(upds)-1].ID + 1
		}
		for i := range upds {
			go lp.handle(b.Handle, &upds[i])
		}
	}
}

func (lp *LongPoller) handle(h Handler, upd *tg.Update) {
	lp.wg.Add(1)
	defer lp.wg.Done()
	if err := h(upd); err != nil {
		lp.cancel()
	}
}
