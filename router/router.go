package router

import (
	"sync"
	"time"

	"github.com/karalef/tgot"
)

// NewRouter makes new initialized queries router.
func NewRouter[Ctx tgot.Context[Ctx], Key comparable, Data any]() *Router[Ctx, Key, Data] {
	return &Router[Ctx, Key, Data]{
		handlers: make(map[Key]handler[Ctx, Key, Data]),
	}
}

// BaseHandler provides base handler methods.
type BaseHandler[Key comparable] interface {
	Name() string

	// Specifies when the handler will be automatically unregistered.
	Timeout() time.Time

	// Called when the handler times out.
	// The current handler will be automatically unregistered so do not
	// call Unregister from this function as this will cause a goroutine leaking.
	Cancel(tgot.BaseContext, Key)
}

// Handler is a handler.
type Handler[Ctx tgot.Context[Ctx], Key comparable, Data any] interface {
	BaseHandler[Key]

	Handle(Ctx, Key, Data)
}

type handler[Ctx tgot.Context[Ctx], Key comparable, Data any] struct {
	Handler[Ctx, Key, Data]
	oneTime bool
}

// Router routes queries by registered keys.
type Router[Ctx tgot.Context[Ctx], Key comparable, Data any] struct {
	handlers map[Key]handler[Ctx, Key, Data]
	mut      sync.Mutex
}

func (r *Router[Ctx, Key, Data]) gc(ctx Ctx) {
	t := time.Now()
	for key, h := range r.handlers {
		if t.Before(h.Timeout()) {
			continue
		}
		r.unregister(key)
		h.Cancel(ctx.WithName(h.Name()), key)
	}
}

// Route routes update.
func (r *Router[Ctx, Key, Data]) Route(ctx Ctx, key Key, data Data) {
	r.mut.Lock()
	r.gc(ctx)
	h, ok := r.handlers[key]
	if !ok {
		r.mut.Unlock()
		return
	}
	if h.oneTime {
		r.unregister(key)
	}
	r.mut.Unlock()

	h.Handle(ctx.WithName(h.Name()), key, data)
}

func (r *Router[Ctx, Key, Data]) reg(key Key, h handler[Ctx, Key, Data]) {
	if h.Handler == nil || h.Timeout().Before(time.Now()) {
		return
	}
	r.mut.Lock()
	r.handlers[key] = h
	r.mut.Unlock()
}

// Register registers handler for key.
func (r *Router[Ctx, Key, Data]) Register(key Key, h Handler[Ctx, Key, Data]) {
	r.reg(key, handler[Ctx, Key, Data]{Handler: h})
}

// RegisterOneTime registers handler for key which will be unregistered after first call.
func (r *Router[Ctx, Key, Data]) RegisterOneTime(key Key, h Handler[Ctx, Key, Data]) {
	r.reg(key, handler[Ctx, Key, Data]{Handler: h, oneTime: true})
}

func (r *Router[Ctx, Key, Data]) unregister(key Key) {
	delete(r.handlers, key)
}

// Unregister deletes handler associated with the key.
func (r *Router[Ctx, Key, Data]) Unregister(key Key) {
	r.mut.Lock()
	r.unregister(key)
	r.mut.Unlock()
}
