package router

import (
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// NewCallbacks makes new inited callback router.
func NewCallbacks() *Callbacks {
	return &Callbacks{
		NewRouter[tgot.CallbackContext, tgot.MessageSignature, *tg.CallbackQuery](),
	}
}

// CallbackHandler represents callbacks handler.
type CallbackHandler interface {
	Name() string

	// If unreg is true, handler will be automatically deleted.
	Handle(tgot.Context, *tg.CallbackQuery) (ans tgot.CallbackAnswer, unreg bool, err error)

	// Specifies when the handler will be automatically unreged.
	Timeout() time.Time

	// Called when the handler times out.
	// The current handler will be automatically unreged so do not
	// call Unreg from this function as this will cause a deadlock.
	Close(tgot.Context, tgot.MessageSignature) error
}

// Callbacks routes callback queries.
type Callbacks struct {
	r *Router[tgot.CallbackContext, tgot.MessageSignature, *tg.CallbackQuery]
}

// Route handles callback query.
//
// It can be used as [Handler.OnCallbackQuery].
func (c *Callbacks) Route(qc tgot.QueryContext[tgot.CallbackAnswer], q *tg.CallbackQuery) {
	c.r.Route(qc, tgot.CallbackSignature(q), q)
}

// Reg registers callback handler for message.
func (c *Callbacks) Reg(sig tgot.MessageSignature, h CallbackHandler, async ...bool) {
	if h != nil {
		c.r.Reg(sig, &callbackWrapper{h}, async...)
	}
}

// Unreg deletes handler associated with the key.
func (c *Callbacks) Unreg(sig tgot.MessageSignature) {
	c.r.Unreg(sig)
}

var _ Handler[tgot.CallbackContext, tgot.MessageSignature, *tg.CallbackQuery] = &callbackWrapper{}

type callbackWrapper struct {
	CallbackHandler
}

func (w *callbackWrapper) Handle(qc tgot.CallbackContext, q *tg.CallbackQuery) (bool, error) {
	ans, unreg, err := w.CallbackHandler.Handle(qc.Context, q)
	if err != nil {
		return unreg, err
	}
	if err = qc.Answer(ans); err != nil {
		qc.Logger().Error("'%s' answer error: %s", w.Name(), err.Error())
	}
	return unreg, nil
}
