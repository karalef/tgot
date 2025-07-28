package updates

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Handler represents update handler.
type Handler interface {
	Allowed() []string
	Handle(*tg.Update) Response
}

// Response sends the api call as a response to the update without returning
// anything.
// If it is not nil, the request to the api will be written in response to the
// webhook or will be sent via api call.
type Response interface {
	Method() string
	Data() *api.Data
}
