package updates

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Poller represents any blocking updates poller.
type Poller interface {
	Run(api *api.API, handler Handler, allowed []string) error
	Close()
}

var _ Poller = &LongPoller{}
var _ Poller = &Webhooker{}

// Handler represents handler function type.
type Handler func(*tg.Update)

// WebhookPoller represents a webhook server that can send api requests in response to webhook requests.
type WebhookPoller interface {
	Poller
	RunWH(api *api.API, handler WHHandler, allowed []string) error
}

var _ WebhookPoller = &Webhooker{}

// WHHandler represents webhook handler function type.
// If the method is not empty, the request to the api will be written in response to the webhook.
type WHHandler func(*tg.Update) (method string, data api.Data)

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
