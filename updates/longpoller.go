package updates

import (
	"context"
	"sync"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

// NewLongPoller creates a LongPoller instance.
func NewLongPoller(filter Filter, timeout, limit, offset int) *LongPoller {
	lp := LongPoller{
		timeout: 30,
		limit:   limit,
		offset:  offset,
		Filter:  filter,
	}
	if timeout > 0 {
		lp.timeout = timeout
	}
	return &lp
}

// LongPoller polls the server for updates via the getUpdates method.
// It must be created via NewLongPoller otherwise only for testing purposes.
type LongPoller struct {
	offset  int
	timeout int
	limit   int

	wg     sync.WaitGroup
	cancel context.CancelFunc

	Filter Filter
}

// Shutdown stops the long poller and waits for all handlers to be completed.
func (lp *LongPoller) Shutdown() {
	lp.Close()
	lp.wg.Wait()
}

// Close stops the long poller.
// It panics if the poller is not running.
func (lp *LongPoller) Close() {
	lp.cancel()
}

// Run starts long polling.
func (lp *LongPoller) Run(a *api.API, h Handler, allowed []string) error {
	if h == nil {
		panic("LongPoller: nil handler")
	}

	ctx := context.Background()
	ctx, lp.cancel = context.WithCancel(ctx)
	defer lp.wg.Wait()

	d := api.NewData()
	d.SetInt("limit", lp.limit)
	d.SetInt("timeout", lp.timeout)
	d.SetJSON("allowed", allowed)

	for {
		d.SetInt("offset", lp.offset)
		upds, err := api.RequestContext[[]tg.Update](ctx, a, "getUpdates", d)
		switch err {
		case nil:
		case context.Canceled, context.DeadlineExceeded:
			return nil
		default:
			return err
		}
		if len(upds) > 0 {
			lp.offset = upds[len(upds)-1].ID
		}
		for i := range filter(upds, lp.Filter) {
			go lp.handle(h, &upds[i])
		}
	}
}

func (lp *LongPoller) handle(h Handler, upd *tg.Update) {
	lp.wg.Add(1)
	h(upd)
	lp.wg.Done()
}
