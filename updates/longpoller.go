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
// timeout must be >0.
// limit must be in the range of 1-100.
func NewLongPoller(timeout, limit, offset int) *LongPoller {
	if timeout < 1 {
		timeout = 30
	}
	if limit < 0 || limit > 100 {
		limit = 0
	}
	lp := LongPoller{
		timeout: timeout,
		limit:   limit,
		offset:  offset,
	}
	return &lp
}

// StartLongPolling creates and runs long poller.
func StartLongPolling(b *tgot.Bot, timeout, limit, offset int) error {
	return NewLongPoller(timeout, limit, offset).Run(b)
}

// LongPoller represents complete Poller that polls the server for updates via the getUpdates method.
type LongPoller struct {
	timeout int
	limit   int
	offset  int

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// Close stops the long poller and waits for all active handlers to complete.
// It panics if the poller is not running.
func (lp *LongPoller) Close() {
	lp.cancel()
	lp.wg.Wait()
}

// Run starts long polling with background context.
func (lp *LongPoller) Run(b *tgot.Bot) error {
	return lp.RunContext(context.Background(), b)
}

// RunContext starts long polling.
func (lp *LongPoller) RunContext(ctx context.Context, b *tgot.Bot) error {
	if b == nil {
		panic("LongPoller: nil bot")
	}

	ctx, lp.cancel = context.WithCancel(ctx)
	defer lp.wg.Wait()

	d := api.NewData()
	d.SetInt("limit", lp.limit)
	d.SetInt("timeout", lp.timeout)
	d.SetJSON("allowed", b.Allowed())

	for a := b.API(); ; {
		d.SetInt("offset", lp.offset)
		upds, err := api.RequestContext[[]tg.Update](ctx, a, "getUpdates", d)
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
