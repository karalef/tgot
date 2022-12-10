package updates

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Poller represents any blocking updates poller.
type Poller interface {
	// should block the goroutine
	Run(api *api.API, handler Handler, allowed []string) error
	Close()
}

// Handler represents handler function type.
type Handler func(*tg.Update)

// FilterFunc represents filter function type.
type FilterFunc func(*tg.Update) bool

// Filter returns only those items to which keep returned true.
func Filter(list []tg.Update, keep FilterFunc) []tg.Update {
	if len(list) == 0 || keep == nil {
		return list
	}
	filtered := list[:0]
	for i := range list {
		if keep(&list[i]) {
			filtered = append(filtered, list[i])
		}
	}
	return filtered
}
