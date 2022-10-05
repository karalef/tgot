package updates

import (
	"context"

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

	Filter func(*tg.Update) bool
}

// Run starts long polling.
func (lp *LongPoller) Run(ctx context.Context, a *api.API, h Handler, allowed []string) error {
	if h == nil {
		panic("LongPoller: nil handler")
	}
	d := api.Data{}
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
			go h(&upds[i])
		}
	}
}
