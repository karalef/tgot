package router

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// NewCallbacks makes new initialized callback router.
func NewCallbacks() *Callbacks {
	return &Callbacks{
		NewRouter[tgot.Query[tgot.CallbackAnswer], tgot.MessageID, *tg.CallbackQuery](),
	}
}

// CallbackHandler represents callbacks handler.
type CallbackHandler interface {
	BaseHandler[tgot.MessageID]

	Handle(*tgot.Message, *tg.CallbackQuery) tgot.CallbackAnswer

	// Called if the tgot.CallbackContext.Answer returns an error.
	OnError(ctx tgot.Query[tgot.CallbackAnswer], err error)
}

// Callbacks routes callback queries.
type Callbacks struct {
	r *Router[tgot.Query[tgot.CallbackAnswer], tgot.MessageID, *tg.CallbackQuery]
}

// Route handles callback query.
//
// It can be used as [Handler.OnCallbackQuery].
func (c *Callbacks) Route(qc tgot.Query[tgot.CallbackAnswer], q *tg.CallbackQuery) {
	c.r.Route(qc, tgot.CallbackMsgID(q), q)
}

// Register registers callback handler for message.
func (c *Callbacks) Register(sig tgot.MessageID, h CallbackHandler) {
	if h != nil {
		c.r.Register(sig, &callbackWrapper{h})
	}
}

// RegisterOneTime registers callback handler for message which will be unregistered after first call.
func (c *Callbacks) RegisterOneTime(sig tgot.MessageID, h CallbackHandler) {
	if h != nil {
		c.r.RegisterOneTime(sig, &callbackWrapper{h})
	}
}

// Unregister deletes handler associated with the key.
func (c *Callbacks) Unregister(sig tgot.MessageID) {
	c.r.Unregister(sig)
}

var _ Handler[tgot.Query[tgot.CallbackAnswer], tgot.MessageID, *tg.CallbackQuery] = &callbackWrapper{}

type callbackWrapper struct {
	CallbackHandler
}

func (w *callbackWrapper) Handle(qc tgot.Query[tgot.CallbackAnswer], sig tgot.MessageID, q *tg.CallbackQuery) {
	ans := w.CallbackHandler.Handle(tgot.WithMessage(qc, tgot.CallbackMsgID(q)), q)
	if err := qc.Answer(ans); err != nil {
		w.CallbackHandler.OnError(qc, err)
	}
}
