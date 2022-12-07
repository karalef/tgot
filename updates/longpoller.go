package updates

import (
	"context"
	"sync"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// NewLongPoller creates a LongPoller instance.
//
// timeout must be >1.
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

// LongPoller polls the server for updates via the getUpdates method.
// It must be created via NewLongPoller otherwise only for testing purposes.
type LongPoller struct {
	offset  int
	timeout int
	limit   int

	wg     sync.WaitGroup
	cancel context.CancelFunc

	Filter FilterFunc
}

// Close stops the long poller and waits for all active handlers to complete.
// It panics if the poller is not running.
func (lp *LongPoller) Close() {
	lp.cancel()
	lp.wg.Wait()
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
		d.SetInt("offset", lp.offset+1)
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
		for i := range Filter(upds, lp.Filter) {
			go lp.handle(h, &upds[i])
		}
	}
}

func (lp *LongPoller) handle(h Handler, upd *tg.Update) {
	lp.wg.Add(1)
	defer lp.wg.Done()
	h(upd)
}
