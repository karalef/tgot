package router

import (
	"sync"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/logger"
)

// NewRouter makes new inited queries router.
func NewRouter[Ctx ctxType[Ctx], Key comparable, Data any]() *Router[Ctx, Key, Data] {
	return &Router[Ctx, Key, Data]{
		handlers: make(map[Key]*handler[Ctx, Key, Data]),
	}
}

type baseContext interface {
	Ctx() tgot.Context
	Logger() *logger.Logger
}

type ctxType[c baseContext] interface {
	baseContext
	Child(string) c
}

// Handler represents handler with timeout.
type Handler[Ctx ctxType[Ctx], Key comparable, Data any] interface {
	Name() string

	// If unreg is true, handler will be automatically deleted.
	Handle(Ctx, Data) (unreg bool, err error)

	// Specifies when the handler will be automatically unreged.
	Timeout() time.Time

	// Called when the handler times out.
	// The current handler will be automatically unreged so do not
	// call Unreg from this function as this will cause a deadlock.
	Close(tgot.Context, Key) error
}

type handler[Ctx ctxType[Ctx], Key comparable, Data any] struct {
	mut     sync.Mutex
	async   bool
	invalid bool
	Handler[Ctx, Key, Data]
}

func (h *handler[Key, Data, Ctx]) lock() {
	h.mut.Lock()
}

func (h *handler[Key, Data, Ctx]) unlock() {
	h.mut.Unlock()
}

// Router routes queries.
type Router[Ctx ctxType[Ctx], Key comparable, Data any] struct {
	handlers map[Key]*handler[Ctx, Key, Data]
	mut      sync.Mutex
}

func (r *Router[Ctx, Key, Data]) gc(ctx tgot.Context) {
	t := time.Now()
	for key, h := range r.handlers {
		if t.Before(h.Timeout()) {
			continue
		}
		delete(r.handlers, key)
		err := h.Close(ctx.Child(h.Name()), key)
		if err != nil {
			ctx.Logger().Error("Close '%s' ended with an error: %s", h.Name(), err.Error())
		}
	}
}

// Route routes update.
func (r *Router[Ctx, Key, Data]) Route(ctx Ctx, key Key, data Data) {
	r.mut.Lock()
	r.gc(ctx.Ctx())
	h, ok := r.handlers[key]
	if ok && !h.async {
		h.lock()
		defer h.unlock()
	}
	if !ok || h.invalid {
		r.mut.Unlock()
		return
	}
	r.mut.Unlock()

	var err error
	h.invalid, err = h.Handle(ctx.Child(h.Name()), data)
	if err != nil {
		ctx.Logger().Error("handler '%s' ended with an error: %s", h.Name(), err.Error())
	}
	if h.invalid {
		r.Unreg(key)
	}
}

// Reg registers handler for key.
func (r *Router[Ctx, Key, Data]) Reg(key Key, h Handler[Ctx, Key, Data], async ...bool) {
	if h != nil || h.Timeout().Before(time.Now()) {
		return
	}
	r.mut.Lock()
	r.handlers[key] = &handler[Ctx, Key, Data]{Handler: h, async: len(async) > 0 && async[0]}
	r.mut.Unlock()
}

// Unreg deletes handler associated with the key.
func (r *Router[Ctx, Key, Data]) Unreg(key Key) {
	r.mut.Lock()
	delete(r.handlers, key)
	r.mut.Unlock()
}
