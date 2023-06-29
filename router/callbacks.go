package router

import (
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// NewCallbacks makes new inited callback router.
func NewCallbacks() *Callbacks {
	return &Callbacks{
		NewRouter[tgot.CallbackContext, tgot.MsgSignature, *tg.CallbackQuery](),
	}
}

// CallbackHandler represents callbacks handler.
type CallbackHandler interface {
	Name() string

	// If unreg is true, handler will be automatically deleted.
	Handle(tgot.MessageContext, *tg.CallbackQuery) (ans tgot.CallbackAnswer, unreg bool)

	// Called if the tgot.CallbackContext.Answer returns an error.
	// if forceUnreg is true the handler will be deleted even if the Handle method has returned false.
	OnError(ctx tgot.CallbackContext, err error) (forceUnreg bool)

	// Specifies when the handler will be automatically unreged.
	Timeout() time.Time

	// Called when the handler times out.
	// The current handler will be automatically unreged so do not
	// call Unreg from this function as this will cause a deadlock.
	Cancel(tgot.MessageContext)
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

// Reg registers callback handler for message.
func (c *Callbacks) Reg(sig tgot.MsgSignature, h CallbackHandler) {
	if h != nil {
		c.r.Reg(sig, &callbackWrapper{h})
	}
}

// Unreg deletes handler associated with the key.
func (c *Callbacks) Unreg(sig tgot.MsgSignature) {
	c.r.Unreg(sig)
}

var _ Handler[tgot.CallbackContext, tgot.MsgSignature, *tg.CallbackQuery] = &callbackWrapper{}

type callbackWrapper struct {
	CallbackHandler
}

func (w *callbackWrapper) Handle(qc tgot.CallbackContext, sig tgot.MsgSignature, q *tg.CallbackQuery) bool {
	ans, unreg := w.CallbackHandler.Handle(qc.OpenMessage(sig), q)
	if err := qc.Answer(ans); err != nil {
		unreg = unreg || w.CallbackHandler.OnError(qc, err)
	}
	return unreg
}

func (w *callbackWrapper) Cancel(ctx tgot.Context, sig tgot.MsgSignature) {
	w.CallbackHandler.Cancel(ctx.OpenMessage(sig))
}
