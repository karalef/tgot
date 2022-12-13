package updates

import "github.com/karalef/tgot/api/tg"

// Handler represents handler function type.
type Handler func(*tg.Update) error

// Filter wraps handler with filter function.
func Filter(handler Handler, filter func(*tg.Update) (keep bool)) Handler {
	if filter == nil {
		return handler
	}
	return func(u *tg.Update) error {
		if filter(u) {
			return handler(u)
		}
		return nil
	}
}
