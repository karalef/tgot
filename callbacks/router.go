package callbacks

import (
	"tghwbot/bot"
	"tghwbot/bot/tg"
)

// NewRouter makes new inited callback router.
func NewRouter() Router {
	return Router{
		msgs: make(map[bot.MessageSignature]handler),
	}
}

// OnCallback is handler function type.
// If unreg is true, handler will be automatically deleted.
type OnCallback func(bot.Context, *tg.CallbackQuery) (ans bot.CallbackAnswer, unreg bool, err error)

type handler struct {
	name string
	f    OnCallback
}

// Router routes callback queries.
type Router struct {
	msgs map[bot.MessageSignature]handler
}

// Route handles callback query.
//
// It can be used as [Handler.OnCallbackQuery].
//
// myHandler.OnCallbackQuery = myCallbackSystem.Push
func (r *Router) Route(qc bot.QueryContext[bot.CallbackAnswer], q *tg.CallbackQuery) {
	log := qc.Logger()
	sig := qc.CallbackSignature(q)
	h, ok := r.msgs[sig]
	if !ok {
		log.Warn("received a message that is not associated with a handler")
		qc.Answer(bot.CallbackAnswer{})
		return
	}
	ans, unreg, err := h.f(qc.Context, q)
	if err != nil {
		log.Error("handler '%s' ended with an error: %s", h.name, err.Error())

	}
	err = qc.Answer(ans)
	if err != nil {
		log.Error("'%s' answer error: %s", h.name, err.Error())
	}
	if unreg {
		r.Unreg(sig)
	}
}

// Reg registers callback handler for message.
func (r *Router) Reg(sig bot.MessageSignature, name string, f OnCallback) {
	r.msgs[sig] = handler{
		name: name,
		f:    f,
	}
}

// Unreg deletes handler associated with a message.
func (r *Router) Unreg(sig bot.MessageSignature) {
	delete(r.msgs, sig)
}
