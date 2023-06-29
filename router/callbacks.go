package router

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// NewCallbacks makes new initialized callback router.
func NewCallbacks() *Callbacks {
	return &Callbacks{
		NewRouter[tgot.CallbackContext, tgot.MsgSignature, *tg.CallbackQuery](),
	}
}

// CallbackHandler represents callbacks handler.
type CallbackHandler interface {
	BaseHandler[tgot.MsgSignature]

	Handle(tgot.MessageContext, *tg.CallbackQuery) tgot.CallbackAnswer

	// Called if the tgot.CallbackContext.Answer returns an error.
	OnError(ctx tgot.CallbackContext, err error)
}

// Callbacks routes callback queries.
type Callbacks struct {
	r *Router[tgot.CallbackContext, tgot.MsgSignature, *tg.CallbackQuery]
}

// Route handles callback query.
//
// It can be used as [Handler.OnCallbackQuery].
func (c *Callbacks) Route(qc tgot.CallbackContext, q *tg.CallbackQuery) {
	c.r.Route(qc, tgot.CallbackSignature(q), q)
}

// Register registers callback handler for message.
func (c *Callbacks) Register(sig tgot.MsgSignature, h CallbackHandler) {
	if h != nil {
		c.r.Register(sig, &callbackWrapper{h})
	}
}

// RegisterOneTime registers callback handler for message which will be unregistered after first call.
func (c *Callbacks) RegisterOneTime(sig tgot.MsgSignature, h CallbackHandler) {
	if h != nil {
		c.r.RegisterOneTime(sig, &callbackWrapper{h})
	}
}

// Unregister deletes handler associated with the key.
func (c *Callbacks) Unregister(sig tgot.MsgSignature) {
	c.r.Unregister(sig)
}

var _ Handler[tgot.CallbackContext, tgot.MsgSignature, *tg.CallbackQuery] = &callbackWrapper{}

type callbackWrapper struct {
	CallbackHandler
}

func (w *callbackWrapper) Handle(qc tgot.CallbackContext, sig tgot.MsgSignature, q *tg.CallbackQuery) {
	ans := w.CallbackHandler.Handle(qc.OpenMessage(sig), q)
	if err := qc.Answer(ans); err != nil {
		w.CallbackHandler.OnError(qc, err)
	}
}
