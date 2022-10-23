package updates

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
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
