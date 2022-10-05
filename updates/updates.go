package updates

import (
	"context"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

// Poller represents any blocking updates poller.
type Poller interface {
	Run(ctx context.Context, api *api.API, handler Handler, allowed []string) error
}

var _ Poller = &LongPoller{}

// Handler represents handler function type.
type Handler func(*tg.Update)

// Filter represents filter function type.
type Filter func(*tg.Update) bool

func filter(slice []tg.Update, keep Filter) []tg.Update {
	if len(slice) == 0 || keep == nil {
		return slice
	}
	filtered := slice[:0]
	for i := range slice {
		if keep(&slice[i]) {
			filtered = append(filtered, slice[i])
		}
	}
	return filtered
}
