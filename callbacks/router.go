package callbacks

import (
	"sync"
	"time"
	"unsafe"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/tg"
)

// NewRouter makes new inited callback router.
func NewRouter() *Router {
	return &Router{
		handlers: make(map[tgot.MessageSignature]*handler),
	}
}

// Handler contains callback handler with parameters.
type Handler interface {
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

type handler struct {
	mut     sync.Mutex
	invalid bool
	Handler
}

func (h *handler) lock() {
	h.mut.Lock()
}

func (h *handler) unlock() {
	h.mut.Unlock()
}

// Router routes callback queries.
type Router struct {
	handlers map[tgot.MessageSignature]*handler
	mut      sync.Mutex
}

func (r *Router) gc(ctx tgot.Context, current tgot.MessageSignature) (cur bool) {
	t := time.Now()
	for sig, h := range r.handlers {
		if t.After(h.Timeout()) {
			if sig == current {
				cur = true
			}
			delete(r.handlers, sig)
			err := h.Close(ctx.Child(h.Name()), sig)
			if err != nil {
				ctx.Logger().Error("Close '%s' ended with an error: %s", h.Name(), err.Error())
			}
		}
	}
	return
}

// Route handles callback query.
//
// It can be used as [Handler.OnCallbackQuery].
//
// myHandler.OnCallbackQuery = myCallbackSystem.Route.
func (r *Router) Route(qc tgot.QueryContext[tgot.CallbackAnswer], q *tg.CallbackQuery) {
	ans := tgot.CallbackAnswer{}
	sig := tgot.CallbackSignature(q)
	r.mut.Lock()
	if r.gc(qc.Context, sig) {
		qc.Answer(ans)
		r.mut.Unlock()
		return
	}
	h, ok := r.handlers[sig]
	if !ok {
		qc.Answer(ans)
		r.mut.Unlock()
		return
	}
	h.lock()
	defer h.unlock()
	if h.invalid {
		qc.Answer(ans)
		r.mut.Unlock()
		return
	}
	r.mut.Unlock()

	log := qc.Logger()
	var err error
	ans, h.invalid, err = h.Handle(qc.Child(h.Name()), q)
	if err != nil {
		log.Error("handler '%s' ended with an error: %s", h.Name(), err.Error())
	}
	if err = qc.Answer(ans); err != nil {
		log.Error("'%s' answer error: %s", h.Name(), err.Error())
	}
	if h.invalid {
		r.Unreg(sig)
	}
}

func isNil(a any) bool {
	return (*[2]uintptr)(unsafe.Pointer(&a))[1] == 0
}

// Reg registers callback handler for message.
func (r *Router) Reg(sig tgot.MessageSignature, h Handler) {
	if isNil(h) || h.Timeout().Before(time.Now()) {
		return
	}
	r.mut.Lock()
	r.handlers[sig] = &handler{Handler: h}
	r.mut.Unlock()
}

// Unreg deletes handler associated with a message.
func (r *Router) Unreg(sig tgot.MessageSignature) {
	r.mut.Lock()
	delete(r.handlers, sig)
	r.mut.Unlock()
}
